// place files you want to import through the `$lib` alias in this folder.


// const baseRoute = "http://jobboard-back:3000/" 
const baseRoute = "/"

export function Advertisement(id,title, description ,companyLogo,wage,address,zipCode ,city ,workTime){
    this.id = id
    this.title = title
    this.description = description
    this.companyLogo = companyLogo
    this.wage = wage
    this.address = address
    this.zipCode = zipCode
    this.city = city
    this.workTime = workTime
}

export async function getAllAds(){
    const response = await fetch(baseRoute+'advertisements');
    jsonAdvertisements = await response.json();
    advertisements = []
    for(jsonAdvertisement of jsonAdvertisements){
        advertisements.push(new Advertisement(jsonAdvertisement.id,jsonAdvertisement.title,
            jsonAdvertisement.description,jsonAdvertisement.company_id,jsonAdvertisement.wage,
            jsonAdvertisement.address, jsonAdvertisement.zipCode,jsonAdvertisement.city,
            jsonAdvertisement.workTime))
    }
    return advertisements
}

export async function submitApply(message,applicant_ID,advertisement_ID){
    const response = await fetch(baseRoute+'application', {
        method: 'POST',
        body: JSON.stringify({message:message,applicant_ID:applicant_ID,advertisement_ID:advertisement_ID}),
        headers: {
            'Content-Type': 'application/json'
        }});
    return response
}

export function sendCredentials(login,password){
    var token
    fetch(baseRoute+'login', {
        method: 'POST',
        body: JSON.stringify({login:login,password:password}),
        headers: {
            'Content-Type': 'application/json'
        }})
        .then(response =>( response.json())
        .then(data =>( token = data.token)))
        .catch(error =>({ error, token: null }))
    console.log(token)
    return token
}

export function createUser(email,password,name,surname,tel,birthDate){
    let success = false;
    fetch(baseRoute+'user', {
        method: 'POST',
        body: JSON.stringify({email:email,password:password,name:name,surname:surname,tel:tel,birthDate:birthDate}),
        headers: {
            'Content-Type': 'application/json'
        }})
        .then(response =>( response.json()))
        .then(data =>( success = true))
        .catch(error =>({ error }))
    return success
}
