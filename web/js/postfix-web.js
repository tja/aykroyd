var emailListApp = new Vue({
  el: '#email-app',
  data: {
    // Model
    domains: []
  },
  mounted: function() {
    this.update()
  },
  methods: {
    // Load and show current state
    update: function() {
      var app = this

      axios
        .get('/api/domains/')
        .then(function (response) { app.domains = app.decorateDomains(response.data) })
        .catch(function (error) { console.log(error) })
    },

    // Create new forward
    createForward: function(domain, event) {
      var app = this

      axios
        .post('/api/domains/' + domain.name + '/forwards/', domain.create)
        .then(function (response) { app.update() })
        .catch(function (error) { console.log(error) })
    },

    // Update existing forward
    updateForward: function(domain, forward, event) {
      var app = this

      axios
        .put('/api/domains/' + domain.name + '/forwards/' + forward.from + '/', { to: forward.to } )
        .then(function (response) { app.update() })
        .catch(function (error) { console.log(error) })
    },

    // Delete existing forward
    deleteForward: function(domain, forward, e) {
      var app = this

      axios
        .delete('/api/domains/' + domain.name + '/forwards/' + forward.from + '/')
        .then(function (response) { app.update() })
        .catch(function (error) { console.log(error) })
    },

    // Decorate domains structure
    decorateDomains: function(domains) {
      _.forEach(domains, function(domain) {
        // Keep original 'to' values
        _.forEach(domain.forwards, function(forward) { forward.toOriginal = forward.to })

        // Placeholder for new forward
        domain.create = { from: '@' + domain.name, to: '' }
      })

      return domains
    },

    // Check if new forward is valid
    isNewForwardValid: function(domain) {
      // Invalid if 'from' or 'to' are empty
      if (domain.create.from.length == 0) { return false }
      if (domain.create.to.length   == 0) { return false }

      // Invalid if 'from' doesn't end in domain
      var postfix = '@' + domain.name
      if (!_.endsWith(domain.create.from, postfix)) { return false }

      // Invalid if 'from; already exists
      var predicate = [ 'from', domain.create.from ]
      if (_.some(domain.forwards, predicate)) { return false }

      return true
    }
  }
});
