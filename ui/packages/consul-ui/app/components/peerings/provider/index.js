import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { getOwner } from '@ember/application';
import { Tab } from 'consul-ui/components/tab-nav';

export default class PeeringsProvider extends Component {
  @service router;
  @service intl;

  get data() {
    return {
      tabs: this.tabs,
    };
  }

  get tabs() {
    const { peer } = this.args;
    const { router } = this;
    const owner = getOwner(this);

    const { isReceiver, Name: name } = peer;
    let tabs = [
      {
        label: 'Imported Services',
        route: 'dc.peers.edit.imported',
        tooltip: this.intl.t('routes.dc.peers.index.detail.imported.tab-tooltip', { name }),
      },
      {
        label: 'Exported Services',
        route: 'dc.peers.edit.exported',
        tooltip: this.intl.t('routes.dc.peers.index.detail.exported.tab-tooltip', { name }),
      },
    ];

    if (isReceiver) {
      tabs = [...tabs, { label: 'Addresses', route: 'dc.peers.edit.addresses' }];
    }

    return tabs.map((tab) => new Tab({ ...tab, currentRouteName: router.currentRouteName, owner }));
  }
}
