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

    // Create new forward
    createForward: function(domain, event) {
      // todo
      console.log("[" + domain.name + "] Create forward from " + domain.create.from + " to " + domain.create.to)
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
