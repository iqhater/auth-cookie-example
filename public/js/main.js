/* window.onload = () => {
    document.forms.login.onsubmit = (e) => {
        e.preventDefault();

        let login = document.querySelector('#login').value;
        let password = document.querySelector('#password').value;

        let data = {
            login,
            password
        };

        sendLogin(data);
    }
};

function sendLogin(data) {
    console.log(data);

    let formData = new FormData();

    for (let param in data) {
        formData.append(param, data[param]);
    }

    fetch('http://localhost:4450/user', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8',
            },
            body: formData
        })
        .then((response) => {
            return response.body;
        })
        .then((data) => {
            console.log(data);
        })
        .catch((error) => {
            console.error(error);
        })
}; */