import { helper } from '@ember/component/helper';

export function objectEntries([obj = {}] /*, hash*/) {
  return Object.entries(obj);
}

export default helper(objectEntries);
