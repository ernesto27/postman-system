<script>
  import logo from './assets/images/logo-universal.png'
  import { MakeRequest } from '../wailsjs/go/main/App.js'

  let selectedMethod = 'GET'
  //let url = 'https://jsonplaceholder.typicode.com/posts'
  let url = 'http://localhost:8080'
  const methods = ['GET', 'POST', 'PUT', 'DELETE']
  let isLoading = false;
  let responseData = '';
  let responseStatus = '';
  let activeTab = 'params'
  let activeResponseTab = 'response'
  
  // Headers management
  let requestHeaders = [{ key: '', value: '' }]
  let responseHeaders = {}

  function addHeader() {
    requestHeaders = [...requestHeaders, { key: '', value: '' }]
  }

  function removeHeader(index) {
    requestHeaders = requestHeaders.filter((_, i) => i !== index)
  }

  // Query params management
  let queryParams = [{ key: '', value: '' }];

  function addQueryParam() {
    queryParams = [...queryParams, { key: '', value: '' }];
  }

  function removeQueryParam(index) {
    queryParams = queryParams.filter((_, i) => i !== index);
    // Update URL after removing parameter
    const newUrl = buildUrl();
    if (newUrl !== url) {
      url = newUrl;
    }
  }

  function buildUrl() {
    try {
      const urlParts = url.split('?')[0];
      const baseUrl = new URL(urlParts);
      
      baseUrl.search = '';
      
      queryParams.forEach(param => {
        if (param.key.trim() && param.value.trim()) {
          baseUrl.searchParams.set(param.key.trim(), param.value.trim());
        }
      });
      
      return baseUrl.toString();
    } catch (error) {
      return url;
    }
  }

  let requestBody = '';
  let bodyType = 'raw';
  let formData = [{ key: '', value: '' }];

  function addFormDataField() {
    formData = [...formData, { key: '', value: '' }];
  }

  function removeFormDataField(index) {
    formData = formData.filter((_, i) => i !== index);
  }

  function handleSend() {
    isLoading = true;
    responseData = '';
    responseStatus = '';
    const finalUrl = buildUrl();
    
    // Build headers object
    const headers = {}
    requestHeaders.forEach(header => {
      if (header.key.trim() && header.value.trim()) {
        headers[header.key.trim()] = header.value.trim()
      }
    })

    // Convert form data array to object
    const formDataObject = {}
    formData.forEach(item => {
      if (item.key.trim() && item.value.trim()) {
        formDataObject[item.key.trim()] = item.value.trim()
      }
    })
    
    MakeRequest(selectedMethod, finalUrl, headers, requestBody, formDataObject, bodyType)
      .then(response => {
        console.log(response)
        responseData = response.body;
        responseHeaders = response.headers;
        responseStatus = response.status;
      })
      .catch(err => {
        console.error(err);
        responseData = `Error: ${err}`;
      })
      .finally(() => {
        isLoading = false;
      });
  }

  function handleKeyPress(event) {
    if (event.key === 'Enter') {
      handleSend();
    }
  }

  function handleParamChange() {
    const newUrl = buildUrl();
    if (newUrl !== url) {
      url = newUrl;
    }
  }

  function handleTabKey(event) {
    if (event.key === 'Tab') {
      event.preventDefault();
      const start = event.target.selectionStart;
      const end = event.target.selectionEnd;
      const value = event.target.value;
      
      event.target.value = value.substring(0, start) + '\t' + value.substring(end);
      event.target.selectionStart = event.target.selectionEnd = start + 1;
    }
  }
</script>

<main class="min-h-screen bg-gray-900 text-gray-100 text-sm">
  <div class="max-w-7xl mx-auto p-6">
    <!-- Header with Logo and URL Bar -->
    <div class="mb-8">
      <img src={logo} alt="Logo" class="h-8 w-auto mb-6">
      <div class="bg-gray-800 rounded-lg p-4">
        <div class="flex flex-col md:flex-row items-center gap-4">
          <select 
            bind:value={selectedMethod}
            class="bg-gray-800 text-white font-medium px-3 py-1.5 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 border border-gray-700 appearance-none cursor-pointer [&>option]:bg-gray-800 [&>option]:text-white text-sm"
            style="background-image: url('data:image/svg+xml;utf8,<svg fill=white xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 24 24%22><path d=%22M7 10l5 5 5-5z%22/></svg>'); background-repeat: no-repeat; background-position: right 8px center; padding-right: 32px;"
          >
            {#each methods as method}
              <option value={method}>{method}</option>
            {/each}
          </select>
          <input 
            bind:value={url}
            type="text" 
            placeholder="Enter request URL"
            on:keypress={handleKeyPress}
            class="flex-1 bg-gray-700 text-gray-300 px-3 py-1.5 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
          >
          <button 
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-1.5 rounded-lg transition-colors disabled:opacity-50 text-sm"
            on:click={handleSend}
            disabled={isLoading}
          >
            {#if isLoading}
              Loading...
            {:else}
              Send
            {/if}
          </button>
        </div>
      </div>
    </div>

    <!-- Request/Response Container -->
    <div class="grid gap-6">
      <!-- Request Section -->
      <div class="bg-gray-800 rounded-lg p-4">
        <div class="border-b border-gray-700 mb-4">
          <div class="flex gap-4">
            <button 
              class="px-3 py-1.5 text-sm {activeTab === 'params' ? 'text-blue-500 border-b-2 border-blue-500' : 'text-gray-400 hover:text-gray-300'}"
              on:click={() => activeTab = 'params'}
            >Params</button>
            <button 
              class="px-3 py-1.5 text-sm {activeTab === 'headers' ? 'text-blue-500 border-b-2 border-blue-500' : 'text-gray-400 hover:text-gray-300'}"
              on:click={() => activeTab = 'headers'}
            >Headers</button>
            <button 
              class="px-3 py-1.5 text-sm {activeTab === 'body' ? 'text-blue-500 border-b-2 border-blue-500' : 'text-gray-400 hover:text-gray-300'}"
              on:click={() => activeTab = 'body'}
            >Body</button>
          </div>
        </div>
        
        <!-- Request Content Area -->
        <div class="bg-gray-900 rounded-lg p-4">
          {#if activeTab === 'params'}
            <div class="space-y-4">
              {#each queryParams as param, i}
                <div class="grid grid-cols-1 md:grid-cols-12 gap-4 items-end">
                  <div class="col-span-2">
                    <label class="block text-xs text-gray-400 mb-1">Key</label>
                    <input 
                      bind:value={param.key}
                      on:input={handleParamChange}
                      type="text" 
                      class="w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                    >
                  </div>
                  <div class="col-span-2">
                    <label class="block text-xs text-gray-400 mb-1">Value</label>
                    <input 
                      bind:value={param.value}
                      on:input={handleParamChange}
                      type="text" 
                      class="w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                    >
                  </div>
                  <div class="col-span-2 flex gap-2">
                    {#if i === queryParams.length - 1}
                      <button 
                        on:click={addQueryParam}
                        class="px-2.5 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm"
                      >
                        Add
                      </button>
                    {/if}
                    {#if queryParams.length > 1}
                      <button 
                        on:click={() => removeQueryParam(i)}
                        class="px-2.5 py-1.5 bg-red-600 hover:bg-red-700 text-white rounded-md text-sm"
                      >
                        Remove
                      </button>
                    {/if}
                  </div>
                </div>
              {/each}
            </div>
          {:else if activeTab === 'headers'}
            <div class="space-y-4">
              {#each requestHeaders as header, i}
                <div class="grid grid-cols-1 md:grid-cols-12 gap-4 items-end">
                  <div class="col-span-2">
                    <label class="block text-xs text-gray-400 mb-1">Key</label>
                    <input 
                      bind:value={header.key}
                      type="text" 
                      class="w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                    >
                  </div>
                  <div class="col-span-2">
                    <label class="block text-xs text-gray-400 mb-1">Value</label>
                    <input 
                      bind:value={header.value}
                      type="text" 
                      class="w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                    >
                  </div>
                  <div class="col-span-2 flex gap-2">
                    {#if i === requestHeaders.length - 1}
                      <button 
                        on:click={addHeader}
                        class="px-2.5 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm"
                      >
                        Add
                      </button>
                    {/if}
                    {#if requestHeaders.length > 1}
                      <button 
                        on:click={() => removeHeader(i)}
                        class="px-2.5 py-1.5 bg-red-600 hover:bg-red-700 text-white rounded-md text-sm"
                      >
                        Remove
                      </button>
                    {/if}
                  </div>
                </div>
              {/each}
            </div>
          {:else if activeTab === 'body'}
            <div class="space-y-4">
              <div class="flex gap-4 mb-4">
                <label class="flex items-center">
                  <input
                    type="radio"
                    bind:group={bodyType}
                    value="raw"
                    class="mr-2"
                  >
                  Raw
                </label>
                <label class="flex items-center">
                  <input
                    type="radio"
                    bind:group={bodyType}
                    value="form-data"
                    class="mr-2"
                  >
                  Form Data
                </label>
              </div>

              {#if bodyType === 'raw'}
                <textarea
                  bind:value={requestBody}
                  placeholder="Enter request body (JSON, text, etc.)"
                  on:keydown={handleTabKey}
                  class="w-full h-40 bg-gray-700 text-gray-300 px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
                ></textarea>
              {:else if bodyType === 'form-data'}
                <div class="space-y-4">
                  {#each formData as field, i}
                    <div class="grid grid-cols-1 md:grid-cols-12 gap-4 items-end">
                      <div class="col-span-2">
                        <label class="block text-xs text-gray-400 mb-1">Key</label>
                        <input 
                          bind:value={field.key}
                          type="text" 
                          class="w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                        >
                      </div>
                      <div class="col-span-2">
                        <label class="block text-xs text-gray-400 mb-1">Value</label>
                        <input 
                          bind:value={field.value}
                          type="text" 
                          class="w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                        >
                      </div>
                      <div class="col-span-2 flex gap-2">
                        {#if i === formData.length - 1}
                          <button 
                            on:click={addFormDataField}
                            class="px-2.5 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm"
                          >
                            Add
                          </button>
                        {/if}
                        {#if formData.length > 1}
                          <button 
                            on:click={() => removeFormDataField(i)}
                            class="px-2.5 py-1.5 bg-red-600 hover:bg-red-700 text-white rounded-md text-sm"
                          >
                            Remove
                          </button>
                        {/if}
                      </div>
                    </div>
                  {/each}
                </div>
              {/if}
            </div>
          {/if}
        </div>
      </div>

      <!-- Response Section -->
      <div class="bg-gray-800 rounded-lg p-4">
        <!-- Add status display -->
        {#if responseStatus}
          <div class="mb-4 flex items-center gap-2">
            <span class="font-semibold">Status:</span>
            <span class={responseStatus.startsWith('2') ? 'text-green-400' : responseStatus.startsWith('4') || responseStatus.startsWith('5') ? 'text-red-400' : 'text-yellow-400'}>
              {responseStatus}
            </span>
          </div>
        {/if}
        
        <div class="border-b border-gray-700 mb-4">
          <div class="flex gap-4">
            <button 
              class="px-3 py-1.5 text-sm {activeResponseTab === 'response' ? 'text-blue-500 border-b-2 border-blue-500' : 'text-gray-400 hover:text-gray-300'}"
              on:click={() => activeResponseTab = 'response'}
            >Response</button>
            <button 
              class="px-3 py-1.5 text-sm {activeResponseTab === 'headers' ? 'text-blue-500 border-b-2 border-blue-500' : 'text-gray-400 hover:text-gray-300'}"
              on:click={() => activeResponseTab = 'headers'}
            >Headers</button>
          </div>
        </div>
        
        <div class="bg-gray-900 rounded-lg p-4">
          <div class="max-h-[600px] max-w-[1200px] overflow-y-auto overflow-x-auto break-word">
            {#if activeResponseTab === 'response'}
              <pre class="text-xs overflow-x-auto"><code class="language-auto">{isLoading ? 'Loading...' : responseData}</code></pre>
            {:else if activeResponseTab === 'headers'}
              <pre class="text-xs overflow-x-auto"><code class="language-auto">{isLoading ? 'Loading...' : JSON.stringify(responseHeaders, null, 2)}</code></pre>
            {/if}
          </div>
        </div>
      </div>

    </div>

  </div>
  
</main>
