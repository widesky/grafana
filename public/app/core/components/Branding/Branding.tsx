import React, { FC } from 'react';
import { css, cx } from '@emotion/css';
import { useTheme2, styleMixins } from '@grafana/ui';
import { colorManipulator } from '@grafana/data';
import config from 'app/core/config';

export interface BrandComponentProps {
  className?: string;
  children?: JSX.Element | JSX.Element[];
}

const LoginLogo: FC<BrandComponentProps> = ({ className }) => {
  const { loginBoxLogo } = config;
  return <img className={className} src={loginBoxLogo} alt="LoginBoxLogo" />;
};

const LoginBackground: FC<BrandComponentProps> = ({ className, children }) => {
  const theme = useTheme2();
  const { loginBackground } = config;

  const background = css`
    &:before {
      content: '';
      position: absolute;
      left: 0;
      right: 0;
      bottom: 0;
      top: 0;
      /// background: url(public/img/g8_login_${theme.isDark ? 'dark' : 'light'}.svg);
      background: url(${loginBackground});
      background-position: top center;
      background-size: auto;
      background-repeat: no-repeat;

      opacity: 0;
      transition: opacity 3s ease-in-out;

      @media ${styleMixins.mediaUp(theme.v1.breakpoints.md)} {
        background-position: bottom;
        background-size: cover;
      }
    }
  `;

  return <div className={cx(background, className)}>{children}</div>;
};

const MenuLogo: FC<BrandComponentProps> = ({ className }) => {
  const { appSidebarLogo } = config;
  return <img className={className} src={appSidebarLogo} alt="ApplicationLogo" />;
};

const LoginBoxBackground = () => {
  const theme = useTheme2();
  const { loginBoxColour } = config;
  return css`
    ///background: ${colorManipulator.alpha(theme.colors.background.primary, 0.7)};
    background: ${loginBoxColour};
    background-size: cover;
  `;
};

const { browserTabTitle } = config;
export class Branding {
  static LoginLogo = LoginLogo;
  static LoginBackground = LoginBackground;
  static MenuLogo = MenuLogo;
  static LoginBoxBackground = LoginBoxBackground;
  static AppTitle = browserTabTitle;
  static LoginTitle = '';
  static GetLoginSubTitle = (): null | string => {
    return null;
  };
}
