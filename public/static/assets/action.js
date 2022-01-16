function starry_night() {
    axios.get('/api/v1/ran_url', {
        responseType: 'json',
    }).then(function (res) {
        window.location.href = res.data.url;
    }).catch(function (err) {
        console.log(err);
    });
}