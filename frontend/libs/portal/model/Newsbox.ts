// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface NewsboxItem {
  _key?: string;
  _id?: string;
  _rev?: string;
  title: string;
  titleDE: string;
  description: string;
  descriptionDE: string;
  image?: string | null;
  link?: string | null;
  expiry?: string | null;
}

export interface NewsboxItemCreateDto {
  title: string;
  titleDE?: string;
  description: string;
  descriptionDE?: string;
  image?: string | null;
  link?: string | null;
  expiry?: string | null;
}

export interface NewsboxItems {
  items: NewsboxItem[];
}

export default interface Newsbox {
  items: NewsboxItem[];
  toShow: number;
}
