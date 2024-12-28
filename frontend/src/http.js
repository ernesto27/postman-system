class HttpClient {
    constructor(baseUrl = '') {
        this.baseUrl = baseUrl;
    }

    setBaseUrl(url) {
        this.baseUrl = url;
    }

    async handleResponse(response) {
        const contentType = response.headers.get('content-type');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        if (contentType && contentType.includes('application/json')) {
            return await response.json();
        }
        return await response.text();
    }

    async get(url, headers = {}) {
        const response = await fetch(this.baseUrl + url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                ...headers
            }
        });
        return this.handleResponse(response);
    }

    async post(url, data, headers = {}) {
        const response = await fetch(this.baseUrl + url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                ...headers
            },
            body: JSON.stringify(data)
        });
        return this.handleResponse(response);
    }

    async put(url, data, headers = {}) {
        const response = await fetch(this.baseUrl + url, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                ...headers
            },
            body: JSON.stringify(data)
        });
        return this.handleResponse(response);
    }

    async delete(url, headers = {}) {
        const response = await fetch(this.baseUrl + url, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                ...headers
            }
        });
        return this.handleResponse(response);
    }
}

export default HttpClient;
