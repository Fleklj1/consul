export default function(obj) {
  // non-recursive for the moment
  return Object.keys(obj).reduce(function(prev, item, i, arr) {
    if (obj[item] !== null) {
      return { ...prev, ...{ [item]: obj[item] } };
    }
    return prev;
  }, {});
}
