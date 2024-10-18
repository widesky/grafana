import React from 'react';
import { useAsync } from 'react-use';

import { config, getBackendSrv } from '@grafana/runtime';
import { Page } from 'app/core/components/Page/Page';

import { AdminSettingsTable } from './AdminSettingsTable';

export type Settings = { [key: string]: { [key: string]: string } };

function AdminSettings() {
  const { loading, value: settings } = useAsync(() => getBackendSrv().get<Settings>('/api/admin/settings'), []);

  return (
    <Page navId="server-settings">
      <Page.Contents>
        <div className="grafana-info-box span8" style={{ margin: '20px 0 25px 0' }} dangerouslySetInnerHTML={{ __html: config.wideSkyWhitelabeling.adminSettingsMessage }}/>

        {loading && <AdminSettingsTable.Skeleton />}

        {settings && <AdminSettingsTable settings={settings} />}
      </Page.Contents>
    </Page>
  );
}

export default AdminSettings;
