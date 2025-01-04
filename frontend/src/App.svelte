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
  let responseTime = 0;
  let responseSize = 0;
  
  // Headers management
  let requestHeaders = [{ key: '', value: '', enabled: true }]
  let responseHeaders = {}

  function addHeader() {
    requestHeaders = [...requestHeaders, { key: '', value: '', enabled: true }]
  }

  function handleHeaderKeyFocus(index) {
    if (index === requestHeaders.length - 1) {
      addHeader();
    }
  }

  // Query params management
  let queryParams = [{ key: '', value: '', enabled: true }];

  function addQueryParam() {
    queryParams = [...queryParams, { key: '', value: '', enabled: true }];
  }

  function handleParamKeyFocus(index) {
    if (index === queryParams.length - 1) {
      addQueryParam();
    }
  }

  function buildUrl() {
    try {
      const urlParts = url.split('?')[0];
      const baseUrl = new URL(urlParts);
      
      baseUrl.search = '';
      
      queryParams.forEach(param => {
        if (param.enabled && param.key.trim() && param.value.trim()) {
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
  let formData = [{ key: '', value: '', type: 'text' }];

  function addFormDataField() {
    formData = [...formData, { key: '', value: '', type: 'text' }];
  }

  function handleFormDataKeyFocus(index) {
    if (index === formData.length - 1) {
      addFormDataField();
    }
  }

  async function handleSend() {
    isLoading = true;
    responseData = '';
    responseStatus = '';
    const finalUrl = buildUrl();
    
    // Build headers object
    const headers = {}
    requestHeaders.forEach(header => {
      if (header.enabled && header.key.trim() && header.value.trim()) {
        headers[header.key.trim()] = header.value.trim()
      }
    })

    // Handle form data and files
    let formDataObj = {};
    let files = [];
    
    if (bodyType === 'form-data') {
      for (const item of formData) {
        if (item.key.trim()) {
          if (item.type === 'file' && item.value instanceof File) {
            // Convert file to base64
            const base64 = await new Promise((resolve) => {
              const reader = new FileReader();
              reader.onloadend = () => resolve(reader.result.split(',')[1]);
              reader.readAsDataURL(item.value);
            });
            
            files.push({
              fieldName: item.key.trim(),
              fileName: item.value.name,
              content: base64
            });
          } else if (item.type === 'text' && item.value.trim()) {
            formDataObj[item.key.trim()] = item.value.trim();
          }
        }
      }
    }
    
    const requestParams = {
      method: selectedMethod,
      url: finalUrl,
      headers: headers,
      body: bodyType === 'raw' ? requestBody : '',
      formData: formDataObj,
      files: files,
      bodyType: bodyType
    };
    
    MakeRequest(requestParams)
      .then(response => {
        console.log(response)
        responseData = response.body;
        responseHeaders = response.headers;
        responseStatus = response.status;
        responseTime = response.timeInSeconds;
        responseSize = response.responseSize;
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

  function formatSize(bytes) {
    if (bytes < 1024) return bytes + ' B';
    else if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' KB';
    else return (bytes / 1048576).toFixed(1) + ' MB';
  }

  function formatTime(seconds) {
    return (seconds * 1000).toFixed(0) + ' ms';
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
            <div class="space-y-4 max-h-[200px] overflow-y-auto pr-2 pb-2">
              {#each queryParams as param, i}
                <div class="grid grid-cols-1 md:grid-cols-[24px_2fr_2fr] gap-4 items-end">
                  <div class="flex items-end">
                    <div class="h-[24px] flex items-center">
                      <input 
                        type="checkbox"
                        bind:checked={param.enabled}
                        on:change={handleParamChange}
                        class="w-3 h-3 rounded bg-gray-700 border-gray-600 text-blue-600 focus:ring-1 focus:ring-blue-600 focus:ring-offset-1 focus:ring-offset-gray-900 cursor-pointer"
                      >
                    </div>
                  </div>
                  <div>
                    <label class="block text-xs text-gray-400 mb-0.5">Key</label>
                    <input 
                      bind:value={param.key}
                      on:input={handleParamChange}
                      on:focus={() => handleParamKeyFocus(i)}
                      type="text" 
                      class="param-input w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm {!param.enabled ? 'opacity-50' : ''}"
                      disabled={!param.enabled}
                    >
                  </div>
                  <div>
                    <label class="block text-xs text-gray-400 mb-0.5">Value</label>
                    <input 
                      bind:value={param.value}
                      on:input={handleParamChange}
                      type="text" 
                      class="param-input w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm {!param.enabled ? 'opacity-50' : ''}"
                      disabled={!param.enabled}
                    >
                  </div>
                </div>
              {/each}
            </div>
          {:else if activeTab === 'headers'}
            <div class="space-y-4 max-h-[200px] overflow-y-auto pr-2 pb-2">
              {#each requestHeaders as header, i}
                <div class="grid grid-cols-1 md:grid-cols-[24px_2fr_2fr] gap-4 items-end">
                  <div class="flex items-end">
                    <div class="h-[24px] flex items-center">
                      <input 
                        type="checkbox"
                        bind:checked={header.enabled}
    
                        class="w-3 h-3 rounded bg-gray-700 border-gray-600 text-blue-600 focus:ring-1 focus:ring-blue-600 focus:ring-offset-1 focus:ring-offset-gray-900 cursor-pointer"
                      >
                    </div>
                  </div>
                  <div>
                    <label class="block text-xs text-gray-400 mb-0.5">Key</label>
                    <input 
                      bind:value={header.key}
                      on:focus={() => handleHeaderKeyFocus(i)}
                      type="text" 
                      class="header-input w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm {!header.enabled ? 'opacity-50' : ''}"
                      disabled={!header.enabled}
                    >
                  </div>
                  <div>
                    <label class="block text-xs text-gray-400 mb-0.5">Value</label>
                    <input 
                      bind:value={header.value}
                      type="text" 
                      class="header-input w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm {!header.enabled ? 'opacity-50' : ''}"
                      disabled={!header.enabled}
                    >
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
                    <div class="grid grid-cols-1 md:grid-cols-[2fr_2fr_3fr] gap-4 items-end">
                      <div>
                        <label class="block text-xs text-gray-400 mb-1">Key</label>
                        <input 
                          bind:value={field.key}
                          type="text"
                          on:focus={() => handleFormDataKeyFocus(i)}
                          class="w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                        >
                      </div>
                      <div>
                        <label class="block text-xs text-gray-400 mb-1">Type</label>
                        <select
                          bind:value={field.type}
                          class="w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm border border-gray-600 appearance-none cursor-pointer"
                          style="background-image: url('data:image/svg+xml;utf8,<svg fill=white xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 24 24%22><path d=%22M7 10l5 5 5-5z%22/></svg>'); background-repeat: no-repeat; background-position: right 8px center; padding-right: 24px;"
                        >
                          <option value="text" class="bg-gray-700">Text</option>
                          <option value="file" class="bg-gray-700">File</option>
                        </select>
                      </div>
                      <div>
                        <label class="block text-xs text-gray-400 mb-1">Value</label>
                        {#if field.type === 'text'}
                          <input 
                            bind:value={field.value}
                            type="text" 
                            class="w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                          >
                        {:else}
                          <input 
                            type="file"
                            class="block w-full text-sm text-gray-300 
                              file:mr-4 file:py-1 file:px-4
                              file:text-sm file:font-medium
                              file:bg-gray-600 file:text-gray-300
                              file:rounded-md file:border-0
                              file:hover:bg-gray-500
                              focus-within:outline-none focus-within:ring-2 focus-within:ring-blue-500
                              cursor-pointer"
                            on:change={(e) => field.value = e.target.files[0]}
                          >
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
          <div class="mb-4 flex items-center gap-4 text-sm">
            <div class="flex items-center gap-2">
              <span class="font-semibold">Status:</span>
              <span class={responseStatus.startsWith('2') ? 'text-green-400' : responseStatus.startsWith('4') || responseStatus.startsWith('5') ? 'text-red-400' : 'text-yellow-400'}>
                {responseStatus}
              </span>
            </div>
            <div class="flex items-center gap-2">
              <span class="font-semibold">Time:</span>
              <span class="text-gray-300">{formatTime(responseTime)}</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="font-semibold">Size:</span>
              <span class="text-gray-300">{formatSize(responseSize)}</span>
            </div>
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
