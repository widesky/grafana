import React, { FC, ReactElement } from 'react';
import { selectors } from '@grafana/e2e-selectors';

import { FormModel } from './LoginCtrl';
import { Button, Form, WideSkyInput, WideSkyField } from '@grafana/ui';
import { css } from '@emotion/css';
import { PasswordField } from '../PasswordField/PasswordField';
import config from 'app/core/config';

interface Props {
  children: ReactElement;
  onSubmit: (data: FormModel) => void;
  isLoggingIn: boolean;
  passwordHint: string;
  loginHint: string;
}

const wrapperStyles = css`
  width: 100%;
  padding-bottom: 16px;
`;

const {
  loginBoxButtonBgColour,
  loginBoxButtonHoverBgColour,
  loginBoxButtonTextColour,
  loginBoxButtonHoverTextColour,
} = config;

export const submitButton = css`
  background: ${loginBoxButtonBgColour};
  justify-content: center;
  width: 100%;
  &:hover {
    background: ${loginBoxButtonHoverBgColour};
    color: ${loginBoxButtonHoverTextColour};
  }
  // For the <span> element to inherit its colour
  color: ${loginBoxButtonTextColour};
`;

export const LoginForm: FC<Props> = ({ children, onSubmit, isLoggingIn, passwordHint, loginHint }) => {
  const { loginBoxTextColour } = config;
  return (
    <div className={wrapperStyles}>
      <Form onSubmit={onSubmit} validateOn="onChange">
        {({ register, errors }) => (
          <>
            <WideSkyField
              label="Email or username"
              invalid={!!errors.user}
              error={errors.user?.message}
              textColour={loginBoxTextColour}
            >
              <WideSkyInput
                {...register('user', { required: 'Email or username is required' })}
                autoFocus
                autoCapitalize="none"
                placeholder={loginHint}
                aria-label={selectors.pages.Login.username}
              />
            </WideSkyField>
            <WideSkyField
              label="Password"
              textColour={loginBoxTextColour}
              invalid={!!errors.password}
              error={errors.password?.message}
            >
              <PasswordField
                id="current-password"
                autoComplete="current-password"
                passwordHint={passwordHint}
                {...register('password', { required: 'Password is required' })}
              />
            </WideSkyField>
            <Button aria-label={selectors.pages.Login.submit} className={submitButton} disabled={isLoggingIn}>
              {isLoggingIn ? 'Logging in...' : 'Log in'}
            </Button>
            {children}
          </>
        )}
      </Form>
    </div>
  );
};
