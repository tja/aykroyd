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
    update: function() {
      var app = this

      // Fetch from server
      axios
        .get('/api/domains/')
        .then(function (response) { app.domains = decorateDomains(response.data) })
        .catch(function (error) { console.log(error) })
    }
  }
});

function decorateDomains(domains) {
  for (var i = 0; i < domains.length; i++) {
    domains[i]['create'] = { from: '', to: '' }
  }

  return domains
}
