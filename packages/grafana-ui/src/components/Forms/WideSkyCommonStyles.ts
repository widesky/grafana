import { css } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';
import { ComponentSize } from '../../types/size';

export const focusCss = () => `
  outline: 2px dotted transparent;
  outline-offset: 2px;
  box-shadow: 0 0 0 0.7px #04275f, 0 0 0px 0.7px $04275f;
  transition: all 0.2s cubic-bezier(0.19, 1, 0.22, 1);
`;

export const getFocusStyle = () => css`
  &:focus {
    ${focusCss()}
  }
`;

export const sharedInputStyle = (theme: GrafanaTheme2, invalid = false) => {
  const borderColor = invalid ? theme.colors.error.border : '#04275F';
  const borderColorHover = invalid ? theme.colors.error.shade : '#04275f';
  const background = '#fff';
  const textColor = '#04275f';

  // Cannot use our normal borders for this color for some reason due the alpha values in them.
  // Need to colors without alpha channel
  const autoFillBorder = '#04275f';

  return css`
    background: ${background};
    line-height: ${theme.typography.body.lineHeight};
    font-size: ${theme.typography.size.md};
    color: ${textColor};
    border: 1px solid ${borderColor};
    padding: ${theme.spacing(0, 1, 0, 1)};

    &:-webkit-autofill,
    &:-webkit-autofill:hover {
      /* Welcome to 2005. This is a HACK to get rid od Chromes default autofill styling */
      box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0), inset 0 0 0 100px ${background}!important;
      -webkit-text-fill-color: ${textColor} !important;
      border-color: ${autoFillBorder};
    }

    &:-webkit-autofill:focus {
      /* Welcome to 2005. This is a HACK to get rid od Chromes default autofill styling */
      box-shadow: 0 0 0 1px ${autoFillBorder}, 0 0 0px 1px ${autoFillBorder}, inset 0 0 0 0.7px rgba(255, 255, 255, 0),
        inset 0 0 0 100px ${background}!important;
      -webkit-text-fill-color: ${textColor} !important;
    }

    &:hover {
      border-color: ${borderColorHover};
    }

    &:disabled {
      background-color: ${theme.colors.action.disabledBackground};
      color: ${theme.colors.action.disabledText};
      border: 1px solid ${theme.colors.action.disabledBackground};

      &:hover {
        border-color: ${borderColor};
      }
    }

    &::placeholder {
      color: ${theme.colors.text.disabled};
      opacity: 1;
    }
  `;
};

export const inputSizes = () => {
  return {
    sm: css`
      width: ${inputSizesPixels('sm')};
    `,
    md: css`
      width: ${inputSizesPixels('md')};
    `,
    lg: css`
      width: ${inputSizesPixels('lg')};
    `,
    auto: css`
      width: ${inputSizesPixels('auto')};
    `,
  };
};

export const inputSizesPixels = (size: string) => {
  switch (size) {
    case 'sm':
      return '200px';
    case 'md':
      return '320px';
    case 'lg':
      return '580px';
    case 'auto':
    default:
      return 'auto';
  }
};

export function getPropertiesForButtonSize(size: ComponentSize, theme: GrafanaTheme2) {
  switch (size) {
    case 'sm':
      return {
        padding: 1,
        fontSize: theme.typography.size.sm,
        height: theme.components.height.sm,
      };

    case 'lg':
      return {
        padding: 3,
        fontSize: theme.typography.size.lg,
        height: theme.components.height.lg,
      };
    case 'md':
    default:
      return {
        padding: 2,
        fontSize: theme.typography.size.md,
        height: theme.components.height.md,
      };
  }
}
