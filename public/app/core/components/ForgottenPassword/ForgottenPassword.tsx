import React, { FC, useState } from 'react';
import {
  Form,
  WideSkyField,
  WideSkyInput,
  Button,
  WideSkyLegend,
  Container,
  useStyles,
  HorizontalGroup,
  LinkButton,
} from '@grafana/ui';
import { getBackendSrv } from '@grafana/runtime';
import { css } from '@emotion/css';
import { GrafanaTheme } from '@grafana/data';
import config from 'app/core/config';

interface EmailDTO {
  userOrEmail: string;
}

const {
  loginBoxButtonBgColour,
  loginBoxButtonHoverBgColour,
  loginBoxButtonTextColour,
  loginBoxButtonHoverTextColour,
} = config;

export const submitButton = css`
  background: ${loginBoxButtonBgColour};
  &:hover {
    background: ${loginBoxButtonHoverBgColour};
    color: ${loginBoxButtonHoverTextColour};
  }
  // For the <span> element to inherit its colour
  color: ${loginBoxButtonTextColour};
`;

const { loginBoxTextColour } = config;
const paragraphStyles = (theme: GrafanaTheme) => css`
  color: ${loginBoxTextColour};
  font-size: ${theme.typography.size.sm};
  font-weight: ${theme.typography.weight.regular};
  margin-top: ${theme.spacing.sm};
  display: block;
`;

const {
  loginBoxLinkButtonBgColour,
  loginBoxLinkButtonHoverBgColour,
  loginBoxLinkButtonTextColour,
  loginBoxLinkButtonHoverTextColour,
} = config;

const backLoginPageStyles = css`
  color: ${loginBoxLinkButtonTextColour};
  background: ${loginBoxLinkButtonBgColour};
  &:hover {
    background: ${loginBoxLinkButtonHoverBgColour};
    color: ${loginBoxLinkButtonHoverTextColour};
  }
`;

export const ForgottenPassword: FC = () => {
  const [emailSent, setEmailSent] = useState(false);
  const styles = useStyles(paragraphStyles);
  const loginHref = `${config.appSubUrl}/login`;
  const { loginBoxTextColour } = config;

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
    <Form onSubmit={sendEmail}>
      {({ register, errors }) => (
        <>
          <WideSkyLegend>Reset password</WideSkyLegend>
          <WideSkyField
            label="User"
            description="Enter your information to get a reset link sent to you"
            invalid={!!errors.userOrEmail}
            error={errors?.userOrEmail?.message}
            textColour={loginBoxTextColour}
          >
            <WideSkyInput
              id="user-input"
              placeholder="Email or username"
              {...register('userOrEmail', { required: 'Email or username is required' })}
            />
          </WideSkyField>
          <HorizontalGroup>
            <Button className={submitButton}>Send reset email</Button>
            <LinkButton className={backLoginPageStyles} fill="text" href={loginHref}>
              Back to login
            </LinkButton>
          </HorizontalGroup>

          <p className={styles}>Did you forget your username or email? Contact your WideSky administrator.</p>
        </>
      )}
    </Form>
  );
};
