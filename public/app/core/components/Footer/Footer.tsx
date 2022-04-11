import React, { FC } from 'react';
import { config } from '@grafana/runtime';
import { Icon, IconName } from '@grafana/ui';

const DEFAULT_GRAFANA_FOOTER_COLOR = '#ccccdc';

export interface FooterLink {
  text: string;
  id?: string;
  icon?: string;
  url?: string;
  target?: string;
  color: string;
}

export let getFooterLinks = (): FooterLink[] => {
  const {
    footerElement1Icon,
    footerElement2Icon,
    footerElement1Text,
    footerElement2Text,
    footerElement1Link,
    footerElement2Link,
    footerElement1TextColor,
    footerElement2TextColor,
  } = config;

  return [
    {
      text: footerElement1Text,
      icon: footerElement1Icon,
      url: footerElement1Link,
      target: '_blank',
      color: footerElement1TextColor,
    },
    {
      text: footerElement2Text,
      icon: footerElement2Icon,
      url: footerElement2Link,
      target: '_blank',
      color: footerElement2TextColor,
    },
  ];
};

export let getVersionLinks = (): FooterLink[] => {
  const {
    buildInfo,
    licenseInfo,
    footerElement3Icon,
    footerElement3Text,
    footerElement3Link,
    footerElement3TextColor,
  } = config;
  const links: FooterLink[] = [];
  const stateInfo = licenseInfo.stateInfo ? ` (${licenseInfo.stateInfo})` : '';
  const footerElement = {
    icon: '',
    text: `${buildInfo.edition}${stateInfo}`,
    url: licenseInfo.licenseUrl,
    color: footerElement3TextColor,
  };

  if (footerElement3Icon !== '') {
    footerElement.icon = footerElement3Icon;
  }

  if (footerElement3Text !== '') {
    footerElement.text = footerElement3Text;
  }

  if (footerElement3Link !== '') {
    footerElement.url = footerElement3Link;
  }

  links.push(footerElement);

  if (buildInfo.hideVersion) {
    return links;
  }

  links.push({ text: `v${buildInfo.version} (${buildInfo.commit})`, color: DEFAULT_GRAFANA_FOOTER_COLOR });

  if (buildInfo.hasUpdate) {
    links.push({
      id: 'updateVersion',
      text: `New version available!`,
      icon: 'download-alt',
      url: 'https://grafana.com/grafana/download?utm_source=grafana_footer',
      target: '_blank',
      color: DEFAULT_GRAFANA_FOOTER_COLOR,
    });
  }

  return links;
};

export function setFooterLinksFn(fn: typeof getFooterLinks) {
  getFooterLinks = fn;
}

export function setVersionLinkFn(fn: typeof getFooterLinks) {
  getVersionLinks = fn;
}

export const Footer: FC = React.memo(() => {
  const links = getFooterLinks().concat(getVersionLinks());
  const { footerPipeColor } = config;

  return (
    <footer className="footer" style={{ color: footerPipeColor }}>
      <div className="text-center">
        <ul>
          {links.map((link) => (
            <li key={link.text}>
              <a href={link.url} target={link.target} rel="noopener">
                {link.icon && <Icon name={link.icon as IconName} />}
                <span style={{ color: link.color }}>{link.text}</span>
              </a>
            </li>
          ))}
        </ul>
      </div>
    </footer>
  );
});

Footer.displayName = 'Footer';
