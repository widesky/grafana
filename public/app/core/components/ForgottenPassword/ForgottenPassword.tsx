import { css } from '@emotion/css';
import React, { useState } from 'react';
import { useForm } from 'react-hook-form';

import { GrafanaTheme2 } from '@grafana/data';
import { getBackendSrv } from '@grafana/runtime';
import { Field, Input, Button, Legend, Container, useStyles2, HorizontalGroup, LinkButton } from '@grafana/ui';
import config from 'app/core/config';

import { getStyles } from '../Login/LoginForm';

interface EmailDTO {
  userOrEmail: string;
}

const paragraphStyles = (theme: GrafanaTheme2) =>
  css({
    color: config.wideSkyWhitelabeling.loginBox.textColor || theme.colors.text.secondary,
    fontSize: theme.typography.bodySmall.fontSize,
    fontWeight: theme.typography.fontWeightRegular,
    marginTop: theme.spacing(1),
    display: 'block',
  });

export const ForgottenPassword = () => {
  const [emailSent, setEmailSent] = useState(false);
  const styles = useStyles2(paragraphStyles);
  const sharedStyles = useStyles2(getStyles);
  const loginHref = `${config.appSubUrl}/login`;
  const {
    handleSubmit,
    register,
    formState: { errors },
  } = useForm<EmailDTO>();

  const sendEmail = async (formModel: EmailDTO) => {
    const res = await getBackendSrv().post('/api/user/password/send-reset-email', formModel);
    if (res) {
      setEmailSent(true);
    }
  };

  if (emailSent) {
    return (
      <div className={styles}>
        <p>An email with a reset link has been sent to the email address. You should receive it shortly.</p>
        <Container margin="md" />
        <LinkButton variant="primary" href={loginHref}>
          Back to login
        </LinkButton>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit(sendEmail)}>
      <Legend className={css({ color: config.wideSkyWhitelabeling.loginBox.textColor })}>Reset password</Legend>
      <Field
        label="User"
        description="Enter your information to get a reset link sent to you"
        invalid={!!errors.userOrEmail}
        error={errors?.userOrEmail?.message}
        className={css({ '& div': { color: config.wideSkyWhitelabeling.loginBox.textColor }, '& span': { color: config.wideSkyWhitelabeling.loginBox.textColor } })}
      >
        <Input
          id="user-input"
          placeholder="Email or username"
          className={sharedStyles.input}
          {...register('userOrEmail', { required: 'Email or username is required' })}
        />
      </Field>
      <HorizontalGroup>
        <Button type="submit" className={sharedStyles.submitButton}>Send reset email</Button>
        <LinkButton className={sharedStyles.linkButton} fill="text" href={loginHref}>
          Back to login
        </LinkButton>
      </HorizontalGroup>

      <p className={styles}>Did you forget your username or email? Contact your {config.wideSkyWhitelabeling.applicationName} administrator.</p>
    </form>
  );
};
