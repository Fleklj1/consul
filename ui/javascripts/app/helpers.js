Ember.Handlebars.helper('panelBar', function(status) {
  var highlightClass;

  if (status == "passing") {
    highlightClass = "bg-green";
  } else {
    highlightClass = "bg-orange";
  }
  return new Handlebars.SafeString('<div class="panel-bar ' + highlightClass + '"></div>');
});

Ember.Handlebars.helper('listBar', function(status) {
  var highlightClass;

  if (status == "passing") {
    highlightClass = "bg-green";
  } else {
    highlightClass = "bg-orange";
  }
  return new Handlebars.SafeString('<div class="list-bar-horizontal ' + highlightClass + '"></div>');
});

Ember.Handlebars.helper('sessionName', function(session) {
  var name;

  if (session.Name === "") {
    name = '<span>' + session.ID + '</span>';
  } else {
    name = '<span>' + session.Name + '</span>' + ' <small>' + session.ID + '</small>';
  }

  return new Handlebars.SafeString(name);
});

Ember.Handlebars.helper('sessionMeta', function(session) {
  var meta = '<div class="metadata">' + session.Behavior + ' behavior</div>';

  if (session.TTL !== "") {
    meta = meta + '<div class="metadata">, ' + session.TTL + ' TTL</div>';
  }

  return new Handlebars.SafeString(meta);
});

Ember.Handlebars.helper('aclName', function(name, id) {
  if (name === "") {
    return id;
  } else {
    return new Handlebars.SafeString(name + ' <small class="pull-right no-case">' + id + '</small>');
  }
});


Ember.Handlebars.helper('formatRules', function(rules) {
  if (rules === "") {
    return "No rules defined";
  } else {
    return rules;
  }
});


// We need to do this because of our global namespace properties. The
// service.Tags
Ember.Handlebars.helper('serviceTagMessage', function(tags) {
  if (tags === null) {
    return "No tags";
  }
});


// Sends a new notification to the UI
function notify(message, ttl) {
  if (window.notifications !== undefined && window.notifications.length > 0) {
    $(window.notifications).each(function(i, v) {
      v.dismiss();
    });
  }
  var notification = new NotificationFx({
    message : '<p>'+ message + '</p>',
    layout : 'growl',
    effect : 'slide',
    type : 'notice',
    ttl: ttl,
  });

  // show the notification
  notification.show();

  // Add the notification to the queue to be closed
  window.notifications = [];
  window.notifications.push(notification);
}

// Tomography

Ember.Handlebars.helper('tomographyGraph', function(tomography, size) {

  // This is ugly, but I'm working around bugs with Handlebars and templating
  // parts of svgs. Basically things render correctly the first time, but when
  // stuff is updated for subsequent go arounds the templated parts don't show.
  // It appears (based on google searches) that the replaced elements aren't
  // being interpreted as http://www.w3.org/2000/svg. Anyway, this works and
  // if/when Handlebars fixes the underlying issues all of this can be cleaned
  // up drastically.

  var max = Math.max.apply(null, tomography.distances);
  var insetSize = size / 2 - 8;
  var buf = '' +
'      <svg width="' + size + '" height="' + size + '">' +
'        <g class="tomography" transform="translate(' + (size / 2) + ', ' + (size / 2) + ')">' +
'          <g>' +
'            <circle class="background" r="' + insetSize + '"/>' +
'            <circle class="axis" r="' + (insetSize * 0.25) + '"/>' +
'            <circle class="axis" r="' + (insetSize * 0.5) + '"/>' +
'            <circle class="axis" r="' + (insetSize * 0.75) + '"/>' +
'            <circle class="border" r="' + insetSize + '"/>' +
'          </g>' +
'          <g class="lines">';
  var sampling = 360 / tomography.n;
  distances = tomography.distances.filter(function () {
    return Math.random() < sampling
  });
  var n = distances.length;
  distances.forEach(function (distance, i) {
    buf += '            <line transform="rotate(' + (i * 360 / n) + ')" y2="' + (-insetSize * (distance / max)) + '"></line>';
  });
  buf += '' +
'          </g>' +
'          <g class="labels">' +
'            <circle class="point" r="5"/>' +
'            <g class="tick" transform="translate(0, ' + (insetSize * -0.25 ) + ')">' +
'              <line x2="70"/>' +
'              <text x="75" y="0" dy=".32em">' + (max > 0 ? (parseInt(max * 25) / 100) : 0) + 'ms</text>' +
'            </g>' +
'            <g class="tick" transform="translate(0, ' + (insetSize * -0.5 ) + ')">' +
'              <line x2="70"/>' +
'              <text x="75" y="0" dy=".32em">' + (max > 0 ? (parseInt(max * 50) / 100) : 0)+ 'ms</text>' +
'            </g>' +
'            <g class="tick" transform="translate(0, ' + (insetSize * -0.75 ) + ')">' +
'              <line x2="70"/>' +
'              <text x="75" y="0" dy=".32em">' + (max > 0 ? (parseInt(max * 75) / 100) : 0) + 'ms</text>' +
'            </g>' +
'            <g class="tick" transform="translate(0, ' + (insetSize * -1) + ')">' +
'              <line x2="70"/>' +
'              <text x="75" y="0" dy=".32em">' + (max > 0 ? (parseInt(max * 100) / 100) : 0) + 'ms</text>' +
'            </g>' +
'          </g>' +
'        </g>' +
'      </svg>';

  return new Handlebars.SafeString(buf);
});
