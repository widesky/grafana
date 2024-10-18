import { css } from '@emotion/css';
import React, { ReactElement, useId } from 'react';
import { useForm } from 'react-hook-form';

import { GrafanaTheme2 } from '@grafana/data';
import { selectors } from '@grafana/e2e-selectors';
import { config } from '@grafana/runtime';
import { Button, Input, Field, useStyles2 } from '@grafana/ui';
import { t } from 'app/core/internationalization';

import { PasswordField } from '../PasswordField/PasswordField';

import { FormModel } from './LoginCtrl';

interface Props {
  children: ReactElement;
  onSubmit: (data: FormModel) => void;
  isLoggingIn: boolean;
  passwordHint: string;
  loginHint: string;
}

export const LoginForm = ({ children, onSubmit, isLoggingIn, passwordHint, loginHint }: Props) => {
  const styles = useStyles2(getStyles);
  const usernameId = useId();
  const passwordId = useId();
  const {
    handleSubmit,
    register,
    formState: { errors },
  } = useForm<FormModel>({ mode: 'onChange' });

  return (
    <div className={styles.wrapper}>
      <form onSubmit={handleSubmit(onSubmit)}>
        <Field
          label={t('login.form.username-label', 'Email or username')}
          invalid={!!errors.user}
          error={errors.user?.message}
          className={css({ '& div': { color: config.wideSkyWhitelabeling.loginBox.textColor } })}
        >
          <Input
            {...register('user', { required: t('login.form.username-required', 'Email or username is required') })}
            id={usernameId}
            autoFocus
            autoCapitalize="none"
            placeholder={loginHint}
            data-testid={selectors.pages.Login.username}
            className={styles.input}
          />
        </Field>
        <Field
          label={t('login.form.password-label', 'Password')}
          invalid={!!errors.password}
          error={errors.password?.message}
          className={css({ '& div': { color: config.wideSkyWhitelabeling.loginBox.textColor } })}
        >
          <PasswordField
            {...register('password', { required: t('login.form.password-required', 'Password is required') })}
            id={passwordId}
            autoComplete="current-password"
            placeholder={passwordHint}
            className={styles.input}
          />
        </Field>
        <Button
          type="submit"
          data-testid={selectors.pages.Login.submit}
          className={styles.submitButton}
          disabled={isLoggingIn}
        >
          {isLoggingIn ? t('login.form.submit-loading-label', 'Logging in...') : t('login.form.submit-label', 'Log in')}
        </Button>
        {children}
      </form>
    </div>
  );
};

export const getStyles = (theme: GrafanaTheme2) => {
  const background = config.wideSkyWhitelabeling.loginBox.textboxBackgroundColor || theme.components.input.background;
  const borderColor = config.wideSkyWhitelabeling.loginBox.textboxBorderColor || theme.colors.background.primary;
  const textColor = config.wideSkyWhitelabeling.loginBox.textColor || theme.components.input.text;

  return {
    wrapper: css({
      width: '100%',
      paddingBottom: theme.spacing(2),
    }),

    submitButton: css({
      justifyContent: 'center',
      width: '100%',
      background: config.wideSkyWhitelabeling.loginBox.buttonBgColor,
      color: config.wideSkyWhitelabeling.loginBox.buttonTextColor,
      '&:hover': {
        background: config.wideSkyWhitelabeling.loginBox.buttonHoverBgColor,
        color: config.wideSkyWhitelabeling.loginBox.buttonHoverTextColor,
      },
    }),

    input: css({
      '& div': {
        '& input': {
          borderColor,
          background,
          color: config.wideSkyWhitelabeling.loginBox.textColor,
          border: `1px solid ${borderColor}`,

          ':hover': {
            background,
            borderColor,
          },

          ':focus': {
            boxShadow: `inset 0 0 0 1px ${borderColor}, inset 0 0 0 100px ${background}!important`,
          },

          '::placeholder': {
            /* Chrome, Firefox, Opera, Safari 10.1+ */
            color: config.wideSkyWhitelabeling.loginBox.textboxPlaceholderColor,
            opacity: 1 /* Firefox */,
          },

          '::-ms-input-placeholder': {
            /* Internet Explorer 10-11 */
            color: config.wideSkyWhitelabeling.loginBox.textboxPlaceholderColor,
          },

          '&:-webkit-autofill, &:-webkit-autofill:hover': {
            /* Welcome to 2005. This is a HACK to get rid od Chromes default autofill styling */
            boxShadow: `inset 0 0 0 1px ${borderColor}, inset 0 0 0 100px ${background}!important`,
            WebkitTextFillColor: `${textColor} !important`,
            borderColor,
          },

          '&:-webkit-autofill:focus': {
            /* Welcome to 2005. This is a HACK to get rid od Chromes default autofill styling */
            boxShadow: `0 0 0 1px ${borderColor}, 0 0 0px 0px rgba(0, 0, 0, 0), inset 0 0 0 1px ${borderColor}, inset 0 0 0 100px ${background}!important`,
            WebkitTextFillColor: `${textColor} !important`,
          },
        },
        '& div': {
          '& button': {
            color: textColor,
          },
        },
      },
    }),

    linkButton: css({
      color: config.wideSkyWhitelabeling.loginBox.linkButtonTextColor,
      background: config.wideSkyWhitelabeling.loginBox.linkButtonBgColor,
      '&:hover': {
        background: config.wideSkyWhitelabeling.loginBox.linkButtonHoverBgColor,
        color: config.wideSkyWhitelabeling.loginBox.linkButtonHoverTextColor,
      },
    }),
  };
};
