function lightning() {
    axios.get('/api/v1/list', {
        responseType: 'json',
    })
        .then(function (res) {
            window.location = res.data[0].url;
        })
        .catch(function (err) {
            console.log(err);
        });

}
