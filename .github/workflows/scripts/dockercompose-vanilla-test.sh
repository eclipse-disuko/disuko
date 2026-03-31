# SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
#
# SPDX-License-Identifier: Apache-2.0

#!/bin/bash
set -euo pipefail

isContainerRunning() {
    local container_name=$1
    health=$(docker inspect --format "{{json .State.Health }}" $container_name | jq -r .Status) # healthy | unhealthy | starting
    status=$(docker inspect --format "{{json .State.Status }}" $container_name | jq -r .) # running | exited
    if [ "$health" = "healthy" ] && [ "$status" = "running" ]; then
        return 0
    else
        return 1
    fi
}

isContainerExitedSuccessful() {
    local container_name=$1
    status=$(docker inspect --format "{{json .State.Status }}" $container_name | jq -r .) # running | exited
    exit_code=$(docker inspect --format "{{.State.ExitCode}}" $container_name) # 0 | 1 | ...
    if [ "$status" = "exited" ] && [ $exit_code -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

docker compose up -d --build >> /dev/null 2>&1
container_names=$(docker ps --format "{{.Names}}" --filter "name=disuko")
remaining_container_count=$(echo "$container_names" | wc -l)
container_checked_count=0
timeout=180
start_time=$(date +%s)
while [ $container_checked_count -lt $remaining_container_count ]; do
echo "$container_checked_count from $remaining_container_count containers are checked..."
for container_name in $container_names; do
  if isContainerRunning "$container_name" || isContainerExitedSuccessful "$container_name"; then
    container_checked_count=$((container_checked_count + 1))
    container_names=$(echo "$container_names" | grep -v $container_name || true)
  fi
done
if [ $(( $(date +%s) - start_time )) -ge $timeout ]; then
  echo "Timeout ${timeout}s reached while waiting for containers to become healthy."
  docker ps
  docker compose down >> /dev/null 2>&1
  exit 1
fi
done
echo "All $remaining_container_count containers are checked."
docker compose down >> /dev/null 2>&1
