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

    // Update existing forward
    updateForward: function(domain, forward, event) {
      // todo
      console.log("[" + domain.name + "] Update forward from " + forward.from + " to " + forward.to)
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
