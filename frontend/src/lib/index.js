//const baseRoute = "http://localhost:3000/"
const baseRoute = "http://"+env.PUBLIC_API_HOST+":"+env.PUBLIC_API_PORT+"/"

export function Advertisement(id, title, description, wage, address, zipCode, city, workTime, companyName, companySiren, companyLogoURL, applied) {
    this.id = id
    this.title = title
    this.description = description
    this.wage = wage
    this.address = address
    this.zipCode = zipCode
    this.city = city
    this.workTime = workTime
    this.companyName = companyName
    this.companySiren = companySiren
    this.companyLogoURL = companyLogoURL
    this.applied = applied
}

export async function getAllAds(token = "", pageCursor = 0, pagePrevious = false) {
    let url = baseRoute + 'advertisements/with_detail?' + new URLSearchParams({
        page_cursor: pageCursor,
        page_previous: pagePrevious,
    })
    let toFetch = fetch(url)
    if (token !== "" && token !== "undefined") {
        toFetch = fetch(url, {
            headers: {
                "Authorization": "Basic " + token
            }
        })
    }
    let promise = toFetch.then(response => (response.json()))
    .then((rep) => {
        let advertisements = []
        rep["Data"].forEach((jsonAdvertisement) => {
            advertisements.push(new Advertisement(
                jsonAdvertisement.ID, jsonAdvertisement.Title,
                jsonAdvertisement.Description, jsonAdvertisement.Wage,
                jsonAdvertisement.Address, jsonAdvertisement.ZipCode, jsonAdvertisement.City,
                jsonAdvertisement.WorkTimeNs / 3600000000000,
                jsonAdvertisement["Company"].Name, jsonAdvertisement["Company"].Siren,
                jsonAdvertisement["Company"].LogoURL, jsonAdvertisement.Applied))
        })
        return [advertisements, rep.Cursors.Previous, rep.Cursors.Next]
    })

return await promise


}

export async function submitApply(message, advertisement_ID, token) {
    let url = baseRoute + 'applications/me?' + new URLSearchParams({
        message: message,
        advertisement_id: advertisement_ID,
    })

    let promise = await fetch(url, {
        method: 'POST',
        headers: {
            "Authorization": "Basic " + token
        },
        mode: "cors",
    })
}

export async function getSvg(fileName) {
    let promise = fetch(fileName)
        .then((res) => res.text())
        .then((text) => {
            return text;
        })
        .catch((e) => console.error(e));
    return await promise;
}

export async function getMe(token) {
    let url = baseRoute + 'users/me'
    let promise = fetch(url, {
        method: 'GET',
        headers: {
            "Authorization": "Basic " + token
        }
    })
        .then(response => (response.json()))
        .then((data) => {
            return data
        })
        .catch((error) => {
            return false
        })
    return await promise
}

export function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) === 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

export async function sendCredentials(login, password) {

    let url = baseRoute + 'users/login?' + new URLSearchParams({
        email: login,
        password: password
    })
    let promise =
        fetch(url, {
            method: 'POST',
        })
            .then(response => (response.json())
                .then(data => {
                    return data.Token
                }))
            .catch(error => {
                console.log(error)
            })
    return await promise
}

export async function createUser(email, password, name, surname, tel, birthDate) {
    let url = baseRoute + 'users?' + new URLSearchParams({
        email: email,
        name: name,
        surname: surname,
        phone: tel,
        date_of_birth_utc: birthDate,
        password: password,
        role: "user"
    })
    let promise = fetch(url, {method: 'POST'})
        .then(response => (response.json())
            .then(data => {
                return data.Token
            }))
        .catch(error => {
            console.log(error)
        })
    return await promise
}

async function updatePwd(password, token) {
    let url = baseRoute + 'users/password/me?' + new URLSearchParams({
        password: password,
    })
    let promise = fetch(url, {
        method: 'PUT',
        headers: {
            "Authorization": "Basic " + token
        }
    })
        .then(data => {
            return true
        })
        .catch(error => {
            return false
        })
    return await promise
}

export async function updateMyInfo(email, name, surname, tel, birthDate, token) {
    let url = baseRoute + 'users/me?' + new URLSearchParams({
        email: email,
        name: name,
        surname: surname,
        phone: tel,
        date_of_birth_utc: birthDate,
        role: "user"
    })
    let promise = fetch(url, {
        method: 'PUT',
        headers: {
            "Authorization": "Basic " + token
        }
    })
        .then(data => {
            return true
        })
        .catch(error => {
            return false
        })
    return await promise
}

export async function updateProfile(email, name, surname, tel, birthDate, token, password) {
    let infoOk = await updateMyInfo(email, name, surname, tel, birthDate, token)
    if (infoOk) {
        return await updatePwd(password, token)
    }
    return false
}