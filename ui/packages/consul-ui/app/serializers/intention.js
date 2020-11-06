import Serializer from './application';
import { inject as service } from '@ember/service';
import { get } from '@ember/object';
import { PRIMARY_KEY, SLUG_KEY } from 'consul-ui/models/intention';
import removeNull from 'consul-ui/utils/remove-null';

export default Serializer.extend({
  primaryKey: PRIMARY_KEY,
  slugKey: SLUG_KEY,
  encoder: service('encoder'),
  init: function() {
    this._super();
    this.uri = this.encoder.uriTag();
  },
  ensureID: function(item) {
    if (!get(item, 'ID.length')) {
      item.Legacy = false;
    } else {
      item.Legacy = true;
      item.LegacyID = item.ID;
    }
    item.ID = this
      .uri`${item.SourceNS}:${item.SourceName}:${item.DestinationNS}:${item.DestinationName}`;
    return item;
  },
  respondForQuery: function(respond, query) {
    return this._super(
      cb =>
        respond((headers, body) => {
          return cb(
            headers,
            body.map(item => this.ensureID(removeNull(item)))
          );
        }),
      query
    );
  },
  respondForQueryRecord: function(respond, query) {
    return this._super(
      cb =>
        respond((headers, body) => {
          body = this.ensureID(removeNull(body));
          return cb(headers, body);
        }),
      query
    );
  },
  respondForCreateRecord: function(respond, serialized, data) {
    const slugKey = this.slugKey;
    const primaryKey = this.primaryKey;
    return respond((headers, body) => {
      body = data;
      body.ID = this
        .uri`${serialized.SourceNS}:${serialized.SourceName}:${serialized.DestinationNS}:${serialized.DestinationName}`;
      return this.fingerprint(primaryKey, slugKey, body.Datacenter)(body);
    });
  },
  respondForUpdateRecord: function(respond, serialized, data) {
    const slugKey = this.slugKey;
    const primaryKey = this.primaryKey;
    return respond((headers, body) => {
      body = data;
      body.LegacyID = body.ID;
      body.ID = serialized.ID;
      return this.fingerprint(primaryKey, slugKey, body.Datacenter)(body);
    });
  },
});
