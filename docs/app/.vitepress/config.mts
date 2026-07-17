import {defineConfig} from 'vitepress';

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: 'DISUKO',
  description: 'Open-source license compliance, without the paperwork',
  base: '/disuko/docs/',
  themeConfig: {
    nav: [
      {text: 'Home', link: '/'},
      {text: 'What is DISUKO?', link: '/introduction'},
      {text: 'Features', link: '/features'},
      {text: 'How it works', link: '/workflow'},
      {text: 'Run it locally', link: '/local-setup'},
    ],

    sidebar: [
      {
        text: 'Overview',
        items: [
          {text: 'What is DISUKO?', link: '/introduction'},
          {text: 'Features', link: '/features'},
          {text: 'How it works', link: '/workflow'},
          {text: 'Glossary', link: '/glossary'},
        ],
      },
      {
        text: 'Getting hands-on',
        items: [{text: 'Run it locally', link: '/local-setup'}],
      },
    ],

    socialLinks: [{icon: 'github', link: 'https://projects.eclipse.org/projects/technology.disuko'}],
  },
});
