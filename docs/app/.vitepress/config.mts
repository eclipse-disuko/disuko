import {defineConfig} from 'vitepress';
import {getChildren} from './util';

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: 'DISUKO',
  description: 'A digital solution replacing manual and paper-based workflows',
  base: '/docs/',
  themeConfig: {
    nav: [
      {text: 'Home', link: '/'},
      {text: 'Quickstart', link: '/misc/quickstart'},
      {text: 'Frontend', link: '/frontend/quick_start.html'},
      {text: 'Backend', link: '/backend/'},
      {text: 'Deployment', link: '/deployment/'},
    ],

    sidebar: [
      {
        text: 'Frontend',
        collapsed: true,
        level: 'deep',
        items: getChildren('.', 'frontend'),
      },
      {
        text: 'Backend',
        collapsed: true,
        level: 'deep',
        items: getChildren('.', 'backend'),
      },
      {
        text: 'Deployment',
        collapsed: true,
        level: 'deep',
        items: getChildren('.', 'deployment'),
      },
      {
        text: 'Misc',
        collapsed: true,
        level: 'deep',
        items: getChildren('.', 'misc'),
      },
    ],

    socialLinks: [{icon: 'github', link: 'https://projects.eclipse.org/projects/technology.disuko'}],
  },
});
