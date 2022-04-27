import React from 'react';
import { connect } from 'react-redux';
import { useAsync } from 'react-use';
import { getBackendSrv } from '@grafana/runtime';
import { NavModel } from '@grafana/data';

import { StoreState } from 'app/types';
import { getNavModel } from 'app/core/selectors/navModel';
import Page from 'app/core/components/Page/Page';

type Settings = { [key: string]: { [key: string]: string } };

interface Props {
  navModel: NavModel;
}

function toWidesky() {
  return (
    <a className="external-link" href="https://widesky.cloud/contact-us/" rel="noreferrer" target="_blank">
      WideSky
    </a>
  );
}

function AdminSettings({ navModel }: Props) {
  const { loading, value: settings } = useAsync(
    () => getBackendSrv().get('/api/admin/settings') as Promise<Settings>,
    []
  );

  return (
    <Page navModel={navModel}>
      <Page.Contents isLoading={loading}>
        <div className="grafana-info-box span8" style={{ margin: '20px 0 25px 0' }}>
          These system settings are managed by the {toWidesky()} team. Contact us to request changes.
        </div>

        {settings && (
          <table className="filter-table">
            <tbody>
              {Object.entries(settings).map(([sectionName, sectionSettings], i) => (
                <React.Fragment key={`section-${i}`}>
                  <tr>
                    <td className="admin-settings-section">{sectionName}</td>
                    <td />
                  </tr>
                  {Object.entries(sectionSettings).map(([settingName, settingValue], j) => (
                    <tr key={`property-${j}`}>
                      <td style={{ paddingLeft: '25px' }}>{settingName}</td>
                      <td style={{ whiteSpace: 'break-spaces' }}>{settingValue}</td>
                    </tr>
                  ))}
                </React.Fragment>
              ))}
            </tbody>
          </table>
        )}
      </Page.Contents>
    </Page>
  );
}

const mapStateToProps = (state: StoreState) => ({
  navModel: getNavModel(state.navIndex, 'server-settings'),
});

export default connect(mapStateToProps)(AdminSettings);
