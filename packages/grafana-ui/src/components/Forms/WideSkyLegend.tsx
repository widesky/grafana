import React, { ReactNode } from 'react';
import { useTheme, stylesFactory } from '../../themes';
import { GrafanaTheme } from '@grafana/data';
import { css, cx } from '@emotion/css';

export interface LabelProps extends React.HTMLAttributes<HTMLLegendElement> {
  children: string | ReactNode;
  description?: string;

  textColour?: string;
}

export const getLegendStyles = stylesFactory((theme: GrafanaTheme, textColour) => {
  return {
    legend: css`
      color: ${textColour};
      font-size: ${theme.typography.heading.h3};
      font-weight: ${theme.typography.weight.regular};
      margin: 0 0 ${theme.spacing.md} 0;
    `,
  };
});

export const WideSkyLegend: React.FC<LabelProps> = ({
  children,
  className,
  textColour = '#04275f',
  ...legendProps
}) => {
  const theme = useTheme();
  const styles = getLegendStyles(theme, textColour);

  return (
    <legend className={cx(styles.legend, className)} {...legendProps}>
      {children}
    </legend>
  );
};
