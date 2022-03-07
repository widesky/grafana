import React, { FC, SyntheticEvent } from 'react';
import { Tooltip, Form, WideSkyField, VerticalGroup, Button } from '@grafana/ui';
import { selectors } from '@grafana/e2e-selectors';
import { submitButton } from '../Login/LoginForm';
import { PasswordField } from '../PasswordField/PasswordField';
interface Props {
  onSubmit: (pw: string) => void;
  onSkip?: (event?: SyntheticEvent) => void;
}
import config from 'app/core/config';
import { css } from '@emotion/css';

interface PasswordDTO {
  newPassword: string;
  confirmNew: string;
}

const {
  loginBoxButtonBgColour,
  loginBoxButtonHoverBgColour,
  loginBoxButtonTextColour,
  loginBoxButtonHoverTextColour,
} = config;

export const skipButton = css`
  background: ${loginBoxButtonBgColour};
  &:hover {
    background: ${loginBoxButtonHoverBgColour};
    color: ${loginBoxButtonHoverTextColour};
  }
  // For the <span> element to inherit its colour
  color: ${loginBoxButtonTextColour};
`;

export const ChangePassword: FC<Props> = ({ onSubmit, onSkip }) => {
  const { loginBoxTextColour } = config;
  const submit = (passwords: PasswordDTO) => {
    onSubmit(passwords.newPassword);
  };
  return (
    <Form onSubmit={submit}>
      {({ errors, register, getValues }) => (
        <>
          <WideSkyField
            label="New password"
            invalid={!!errors.newPassword}
            error={errors?.newPassword?.message}
            textColour={loginBoxTextColour}
          >
            <PasswordField
              id="new-password"
              autoFocus
              autoComplete="new-password"
              {...register('newPassword', { required: 'New Password is required' })}
            />
          </WideSkyField>
          <WideSkyField
            label="Confirm new password"
            invalid={!!errors.confirmNew}
            error={errors?.confirmNew?.message}
            textColour={loginBoxTextColour}
          >
            <PasswordField
              id="confirm-new-password"
              autoComplete="new-password"
              {...register('confirmNew', {
                required: 'Confirmed Password is required',
                validate: (v: string) => v === getValues().newPassword || 'Passwords must match!',
              })}
            />
          </WideSkyField>
          <VerticalGroup>
            <Button type="submit" className={submitButton}>
              Submit
            </Button>

            {onSkip && (
              <Tooltip
                content="If you skip you will be prompted to change password next time you log in."
                placement="bottom"
              >
                <Button
                  className={skipButton}
                  fill="text"
                  onClick={onSkip}
                  type="button"
                  aria-label={selectors.pages.Login.skip}
                >
                  Skip
                </Button>
              </Tooltip>
            )}
          </VerticalGroup>
        </>
      )}
    </Form>
  );
};
