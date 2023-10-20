export async function getSvg(fileName) {
    let promise = fetch("src/lib/Images/" + fileName + '.svg')
        .then((res) => res.text())
        .then((text) => {
            return text;
        })
        .catch((e) => console.error(e));
    return await promise;
}