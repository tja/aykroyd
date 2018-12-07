var emailListApp = new Vue({
  el: '#email-app',
  data: {
    // All data
    domains: []
  },
  mounted: function() {
    // Initial load
    this.update()
  },
  methods: {
    // Load and show current state
    update: function() {
      var app = this

      // Fetch from server
      axios
        .get('/api/domains/')
        .then(function (response) { app.domains = decorateDomains(response.data) })
        .catch(function (error) { console.log(error) })
    },

    // Check if new forwaerd is valid
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
    },

    // Create new forward
    createForward: function(domain, event) {
      var app = this

      // Send to server
      axios
        .post('/api/domains/' + domain.name + '/forwards/', domain.create)
        .then(function (response) { app.update() })
        .catch(function (error) { console.log(error) })
    },

    // Update existing forward
    updateForward: function(domain, forward, event) {
      // todo
      console.log("[" + domain.name + "] Update forward from " + forward.from + " to " + forward.to)
    },

    // Delete existing forward
    deleteForward: function(domain, forward, e) {
      // todo
      console.log("[" + domain.name + "] Delete forward from " + forward.from)
    }
  }
});

function decorateDomains(domains) {
  for (var i = 0; i < domains.length; i++) {
    for (var j = 0; j < domains[i].forwards.length; j++) {
      // Keep original 'to' value
      domains[i].forwards[j].toOriginal = domains[i].forwards[j].to
    }

    // Placeholder for new forward
    domains[i].create = {
      from: '@' + domains[i].name,
      to:   ''
    }
  }

  return domains
}
