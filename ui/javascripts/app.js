window.App = Ember.Application.create({
  rootElement: "#app",
  LOG_TRANSITIONS: true,
  currentPath: ''
});

App.Router.map(function() {

  this.resource("dc", {path: "/:dc"}, function() {
    this.resource("services", { path: "/services" }, function(){
      this.route("show", { path: "/:name" });
    });
    this.resource("nodes", { path: "/nodes" }, function() {
      this.route("show", { path: "/:name" });
    });
  });

  this.route("index", { path: "/" });
  this.route("kv", { path: "/:dc/kv" });
});

//
// The main controller for the application. To use this in other
// controller you have to use the needs api
//
App.ApplicationController = Ember.Controller.extend({
  //
  // Sets the current datacenter to be used when linking. DC is
  // a required parameter in the routes, above, and we need
  // to know our scope. This is so #link-to 'node' dc node-name
  // works.
  //
  setDc: function(dc) {
    localStorage.setItem("current_dc", dc);
  },

  //
  // Retrieves the current datacenter set by this.setDc from
  // localstorage. returns null if the dc has not been set.
  //
  getDc: function() {
    return localStorage.getItem("current_dc");
  }.property("getDc")
});

//
// Superclass to be used by all of the main routes below. All routes
// but the IndexRoute share the need to have a datacenter set.
//
//
App.BaseRoute = Ember.Route.extend({
  //
  // When activating the base route, if we don't have a datacenter set,
  // transition the user to the index route to choose a datacenter.
  //
  activate: function() {
    var controller = this.controllerFor('application');
    if (controller.getDc === null) {
      this.transitionTo('index');
    };
  }
});

//
// The route for choosing datacenters, typically the first route loaded.
//
// Note: This *does not* extend from BaseRoute as that could cause
// and loop of transitions.
//
App.IndexRoute = Ember.Route.extend({
  model: function() {
    return window.fixtures.dcs;
  },
  afterModel: function(dcs, transition) {
    if (dcs.get('length') === 1) {
      this.get('controllers.application').setDc(dcs[0])
      this.transitionTo('services', dcs[0]);
    }
  }
});

//
// path: /
//
// The index is for choosing datacenters.
//
App.IndexController = Ember.Controller.extend({
  needs: ['application'],

  actions: {
    //
    // selectDc is called with the datacenter name to be set for
    // future actions within the application. It's a namespace, essentially.
    //
    // See ApplicationController#setDc
    //
    selectDc: function(dc) {
      this.get('controllers.application').setDc(dc)
      this.transitionToRoute('services', dc)
    }
  }
});

//
// A Consul service.
//
App.Service = Ember.Object.extend({
  //
  // The number of failing checks within the service.
  //
  failingChecks: function() {
    return this.get('Checks').filterBy('Status', 'critical').get('length');
  }.property('failingChecks'),

  //
  // The number of passing checks within the service.
  //
  passingChecks: function() {
    return this.get('Checks').filterBy('Status', 'passing').get('length');
  }.property('passingChecks'),

  //
  // The formatted message returned for the user which represents the
  // number of checks failing or passing. Returns `1 passing` or `2 failing`
  //
  checkMessage: function() {
    if (this.get('hasFailingChecks') === false) {
      return this.get('passingChecks') + ' passing';
    } else {
      return this.get('failingChecks') + ' failing';
    }
  }.property('checkMessage'),

  //
  // Boolean of whether or not there are failing checks in the service.
  // This is used to set color backgrounds and so on.
  //
  hasFailingChecks: function() {
    return (this.get('failingChecks') > 0);
  }.property('hasFailingChecks')
});

//
// A Consul Node
//
App.Node = Ember.Object.extend({
  //
  // The number of failing checks within the service.
  //
  failingChecks: function() {
    return this.get('Checks').filterBy('Status', 'critical').get('length');
  }.property('failingChecks'),

  //
  // The number of passing checks within the service.
  //
  passingChecks: function() {
    return this.get('Checks').filterBy('Status', 'passing').get('length');
  }.property('passingChecks'),

  //
  // The formatted message returned for the user which represents the
  // number of checks failing or passing. Returns `1 passing` or `2 failing`
  //
  checkMessage: function() {
    if (this.get('hasFailingChecks') === false) {
      return this.get('passingChecks') + ' passing';
    } else {
      return this.get('failingChecks') + ' failing';
    }
  }.property('checkMessage'),

  //
  // Boolean of whether or not there are failing checks in the service.
  // This is used to set color backgrounds and so on.
  //
  hasFailingChecks: function() {
    return (this.get('failingChecks') > 0);
  }.property('hasFailingChecks')
});

//
// Display all the services, allow to drill down into the specific services.
//
App.ServicesRoute = App.BaseRoute.extend({
  //
  // Set the services as the routes default model to be called in
  // the template as {{model}}
  //
  setupController: function(controller, model) {
      //
      // Since we have 2 column layout, we need to also display the
      // list of services on the left. Hence setting the attribute
      // {{services}} on the controller.
      //
      controller.set('services', [App.Service.create(window.fixtures.services[0]), App.Service.create(window.fixtures.services[1])]);
  }
});


App.DcRoute = App.BaseRoute.extend({
  //
  // Set the model on the route. We look up the specific service
  // by it's identifier passed via the route
  //
  model: function(params) {
    return params.dc;
  }
});

//
// Display an individual service, as well as the global services in the left
// column.
//
App.ServicesShowRoute = App.BaseRoute.extend({
  //
  // Set the model on the route. We look up the specific service
  // by it's identifier passed via the route
  //
  model: function(params) {
    return [App.Node.create(window.fixtures.services_full[params.name][0]), App.Node.create(window.fixtures.services_full[params.name][1])];
  },

  setupController: function(controller, model) {
      controller.set('content', model);
      //
      // Since we have 2 column layout, we need to also display the
      // list of services on the left. Hence setting the attribute
      // {{services}} on the controller.
      //
      controller.set('services', [App.Service.create(window.fixtures.services[0]), App.Service.create(window.fixtures.services[1])]);
  }
});

//
// DC
//
App.DcView = Ember.View.extend({
    templateName: 'dc',
})

//
// Services
//
App.ServicesView = Ember.View.extend({
    templateName: 'services',
})

//
// Services
//
App.ServicesShowView = Ember.View.extend({
    //
    // We use the same template as we do for the services
    // array and have a simple conditional to display the nested
    // individual service resource.
    //
    templateName: 'service'
})


//
// path: /:dc
//
App.DcController = Ember.Controller.extend({
  needs: ['application']
})

//
//
// path: /:dc/services
//
App.ServicesController = Ember.Controller.extend({
  needs: ['application']
})

//
// path: /:dc/services/:name
//
App.ServicesShowController = Ember.Controller.extend({
  needs: ['application']
})

//=
// path: /:dc/nodes
//
App.NodesController = Ember.Controller.extend({
  needs: ['application']
})

//
// path: /:dc/nodes/:name
//
App.NodesShowController = Ember.Controller.extend({
  needs: ['application']
})

//
// Nodes
//
App.NodesView = Ember.View.extend({
    templateName: 'nodes'
})

App.NodesShowView = Ember.View.extend({
    //
    // We use the same template as we do for the nodes
    // array and have a simple conditional to display the nested
    // individual node resource.
    //
    templateName: 'node'
})


//
// Display an individual node, as well as the global nodes in the left
// column.
//
App.NodesShowRoute = App.BaseRoute.extend({
  //
  // Set the model on the route. We look up the specific node
  // by it's identifier passed via the route
  //
  model: function(params) {
    return App.Node.create(window.fixtures.nodes_full[params.name]);
  },

  setupController: function(controller, model) {
      controller.set('content', model);
      //
      // Since we have 2 column layout, we need to also display the
      // list of nodes on the left. Hence setting the attribute
      // {{nodes}} on the controller.
      //
      controller.set('nodes', [App.Node.create(window.fixtures.nodes[0]), App.Node.create(window.fixtures.nodes[1])]);
  }
});

//
// Display all the nodes, allow to drill down into the specific nodes.
//
App.NodesRoute = App.BaseRoute.extend({
  //
  // Set the node as the routes default model to be called in
  // the template as {{model}}. This is the "expanded" view.
  //
  setupController: function(controller, model) {
      //
      // Since we have 2 column layout, we need to also display the
      // list of nodes on the left. Hence setting the attribute
      // {{nodes}} on the controller.
      //
      controller.set('nodes', [App.Node.create(window.fixtures.nodes[0]), App.Node.create(window.fixtures.nodes[1])]);
  }
});
