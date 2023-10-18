// const baseRoute = "http://jobboard-back:3000/"

const baseRoute = "http://localhost:3000/"

export function Advertisement(id, title, description, wage, address, zipCode, city, workTime, companyName, companySiren, companyLogoURL) {
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
}

export async function getAllAds() {
    let promise = fetch(baseRoute + 'advertisements/with_company')
        .then(response => (response.json()))
        .then((data) => {
            let advertisements = []
            data.forEach((jsonAdvertisement) => {
                advertisements.push(new Advertisement(
                    jsonAdvertisement.ID, jsonAdvertisement.Title,
                    jsonAdvertisement.Description, jsonAdvertisement.Wage,
                    jsonAdvertisement.Address, jsonAdvertisement.ZipCode, jsonAdvertisement.City,
                    jsonAdvertisement.WorkTimeNs / 3600000000000,
                    jsonAdvertisement["Company"].Name, jsonAdvertisement["Company"].Siren,
                    jsonAdvertisement["Company"].LogoURL))
            })
            return advertisements
        })

    return await promise


}

export async function submitApply(message, applicant_ID, advertisement_ID,token) {
    let url = baseRoute + 'applications?' + new URLSearchParams({
        message: message,
        applicant_id: applicant_ID,
        advertisement_id: advertisement_ID,
    })

    let promise = await fetch(url, {
        method: 'POST',
        headers: {
            "Authorization": "Basic "+token
        }
    })
}

export async function getMyId(token) {
    let url = baseRoute + 'usersGetMe'
    let promise = fetch(url, {
        method: 'GET',
        headers: {
            "Authorization": "Basic "+token
        }
    })
        .then(response => (response.json()))
        .then((data) => {
            return data.ID
        })
    return await promise
}
export function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++) {
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
            return data
        }))
        .catch(error => {
            return false
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
        password: password
    })
    let promise = fetch(url, {method: 'POST'})
        .then(data => {
            return true
        })
        .catch(error => {
            return false
        })
    return await promise
}
