import React from 'react';

import { LinkTarget } from '@grafana/data';
import { config } from '@grafana/runtime';
import { Icon, IconName } from '@grafana/ui';
import { t } from 'app/core/internationalization';

export interface FooterLink {
  target: LinkTarget;
  text: string;
  id: string;
  icon?: IconName;
  url?: string;
  color?: string;
  hide?: boolean;
}

export let getFooterLinks = (): FooterLink[] => {
    return config.wideSkyWhitelabeling.footerItems.map(item => ({
        ...item,
        target: '_blank',
        text: t(`nav.help/${item.text.toLocaleLowerCase()}`, item.text),
        id: item.text.toLocaleLowerCase(),
    }));
};

export function getVersionMeta(version: string) {
  const isBeta = version.includes('-beta');

  return {
    hasReleaseNotes: true,
    isBeta,
  };
}

export function getVersionLinks(hideEdition?: boolean): FooterLink[] {
  const { buildInfo, licenseInfo } = config;
  const links: FooterLink[] = [];
  const stateInfo = licenseInfo.stateInfo ? ` (${licenseInfo.stateInfo})` : '';

  if (!hideEdition) {
    links.push({
      target: '_blank',
      id: 'license',
      text: `${buildInfo.edition}${stateInfo}`,
      url: licenseInfo.licenseUrl,
    });
  }

  if (buildInfo.hideVersion) {
    return links;
  }

  const { hasReleaseNotes } = getVersionMeta(buildInfo.version);

  links.push({
    target: '_blank',
    id: 'version',
    text: buildInfo.versionString,
    url: hasReleaseNotes ? `https://github.com/grafana/grafana/blob/main/CHANGELOG.md` : undefined,
  });

  if (buildInfo.hasUpdate) {
    links.push({
      target: '_blank',
      id: 'updateVersion',
      text: `New version available!`,
      icon: 'download-alt',
      url: 'https://grafana.com/grafana/download?utm_source=grafana_footer',
    });
  }

  return links;
}

export function setFooterLinksFn(fn: typeof getFooterLinks) {
  getFooterLinks = fn;
}

export interface Props {
  /** Link overrides to show specific links in the UI */
  customLinks?: FooterLink[] | null;
  hideEdition?: boolean;
}

export const Footer = React.memo(({ customLinks, hideEdition }: Props) => {
  const links = (customLinks || getFooterLinks()).concat(getVersionLinks(hideEdition));

  return (
    <footer className="footer">
      <div className="text-center">
        <ul>
          {links.filter(link => link.hide !== true || link.hide === undefined).map((link, index) => (
            <li key={index} style={{ color: config.wideSkyWhitelabeling.footerPipeColor }}>
              <FooterItem item={link} />
            </li>
          ))}
        </ul>
      </div>
    </footer>
  );
});

Footer.displayName = 'Footer';

function FooterItem({ item }: { item: FooterLink }) {
  const content = item.url ? (
    <a href={item.url} target={item.target} rel="noopener noreferrer" id={item.id}>
      <span style={{ color: item.color }}>{item.text}</span>
    </a>
  ) : (
    <span style={{ color: item.color }}>{item.text}</span>
  );

  return (
    <>
      {item.icon && <Icon name={item.icon} />} {content}
    </>
  );
}
