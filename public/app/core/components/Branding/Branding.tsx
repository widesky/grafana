import { css, cx } from '@emotion/css';
import React, { FC } from 'react';

import { colorManipulator } from '@grafana/data';
import { config } from '@grafana/runtime';
import { useTheme2 } from '@grafana/ui';

export interface BrandComponentProps {
  className?: string;
  children?: JSX.Element | JSX.Element[];
}

export const LoginLogo: FC<BrandComponentProps & { logo?: string }> = ({ className, logo }) => {
  return <img className={className} src={config.wideSkyWhitelabeling.loginBox.logo} alt={config.wideSkyWhitelabeling.applicationName} />;
};

const LoginBackground: FC<BrandComponentProps> = ({ className, children }) => {
  const theme = useTheme2();

  const background = css({
    '&:before': {
      content: '""',
      position: 'fixed',
      left: 0,
      right: 0,
      bottom: 0,
      top: 0,
      background: `url(${config.wideSkyWhitelabeling.loginBackground.image})`,
      backgroundPosition: 'top center',
      backgroundSize: 'auto',
      backgroundRepeat: 'no-repeat',

      opacity: 0,
      transition: 'opacity 3s ease-in-out',

      [theme.breakpoints.up('md')]: {
        backgroundPosition: config.wideSkyWhitelabeling.loginBackground.position,
        backgroundSize: 'cover',
      },
    },
  });

  return <div className={cx(background, className)}>{children}</div>;
};

const MenuLogo: FC<BrandComponentProps> = ({ className }) => {
  return <img className={className} src={config.wideSkyWhitelabeling.appSidebarLogo} alt={config.wideSkyWhitelabeling.applicationName} />;
};

const LoginBoxBackground = () => {
  const theme = useTheme2();
  return css({
    background: config.wideSkyWhitelabeling.loginBox.color,
    backgroundSize: 'cover',
  });
};

export class Branding {
  static LoginLogo = LoginLogo;
  static LoginBackground = LoginBackground;
  static MenuLogo = MenuLogo;
  static LoginBoxBackground = LoginBoxBackground;
  static AppTitle = config.wideSkyWhitelabeling.browserTabTitle;
  static LoginTitle = config.wideSkyWhitelabeling.browserTabSubtitle;
  static HideEdition = config.buildInfo.hideVersion;
  static GetLoginSubTitle = (): null | string => {
    return null;
  };
}
