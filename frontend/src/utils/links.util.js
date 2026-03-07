import axios from 'axios'

export async function getWebsiteMeta(url) {
    const { data } = await axios.get(url);
    const parser = new DOMParser();
    const doc = parser.parseFromString(data, "text/html");

    const title = doc.querySelector("title")?.textContent?.trim() || null;
    const description =
        doc.querySelector('meta[name="description"]')?.getAttribute("content")?.trim() ||
        doc.querySelector('meta[property="og:description"]')?.getAttribute("content")?.trim() ||
        null;

    return [title, description];
}