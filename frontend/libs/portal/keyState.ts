export let CONTROL_IS_PRESSED = false;
export let SHIFT_IS_PRESSED = false;

export function releaseKeys() {
  CONTROL_IS_PRESSED = false;
  SHIFT_IS_PRESSED = false;
}

window.addEventListener('blur', () => releaseKeys, false);

window.addEventListener('focusout', () => releaseKeys, false);

window.addEventListener('onmouseleave', () => releaseKeys, false);

window.addEventListener(
  'onmouseover',
  (e) => {
    releaseKeys();
    const e2 = e as MouseEvent;
    if (e2.shiftKey) {
      SHIFT_IS_PRESSED = true;
    }
    if (e2.ctrlKey) {
      CONTROL_IS_PRESSED = true;
    }
  },
  false,
);

window.addEventListener(
  'keydown',
  (e) => {
    if (!e.key) {
      return;
    }
    switch (e.key.toLowerCase()) {
      case 'control':
        CONTROL_IS_PRESSED = true;
        break;
      case 'shift':
        SHIFT_IS_PRESSED = true;
        break;
    }
  },
  false,
);

window.addEventListener(
  'keyup',
  (e) => {
    if (!e.key) {
      return;
    }
    switch (e.key.toLowerCase()) {
      case 'control':
        CONTROL_IS_PRESSED = false;
        break;
      case 'shift':
        SHIFT_IS_PRESSED = false;
        break;
    }
  },
  false,
);
