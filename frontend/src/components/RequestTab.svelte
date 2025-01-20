<script lang="ts">
  import { onMount } from 'svelte'
  import { MakeRequest } from '../../wailsjs/go/main/App.js'
  import type { PostmanCollection, PostmanRequest, PostmanItem } from '../types/postman.ts';


  const defaultCollection: PostmanCollection = {
    id: "1",
    name: "API Testing Collection",
    items: [
      {
        id: "req1",
        name: "Get Users",
        request: {
          method: "GET",
          url: {
            raw: "https://api.example.com/users",
            protocol: "https",
            host: ["api", "example", "com"],
            path: ["users"],
            port: "8080",
            query: [
              {
                key: "",
                value: "",
                disabled: false
              }
            ]
          },
          header: [
            {
              key: "",
              value: "",
              disabled: false
            },
          ],
          body: {
            mode: "raw",
            raw: "",
            formdata: [{
              key: "",
              value: "",
              type: "text"
            }]
          },
        },
        response: []
      },
    ],
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  };



  export let tabID: string
  export let savedState: PostmanItem | null = null
  export let onStateChange = (state) => {
    console.log(state)
  }

  const methods = ['GET', 'POST', 'PUT', 'DELETE']
  let isLoading = false
  let responseData = ''
  let responseStatus = ''
  let activeTab = 'params'
  let activeResponseTab = 'response'
  let responseTime = 0
  let responseSize = 0
  let responseHeaders = {}


  let item: PostmanItem | null = {
    id: "default",
    name: "Create New User",
    request: {
      method: "POST",
      url: {
        raw: "https://api.example.com/users?active=true",
        protocol: "https",
        host: ["api", "example", "com"],
        path: ["users"],
        port: "8080",
        query: [
          {
            key: "active",
            value: "true",
            disabled: false
          }
        ]
      },
      header: [
        {
          key: "Content-Type",
          value: "application/json",
          disabled: false
        },
        {
          key: "Authorization",
          value: "Bearer xyz123",
          disabled: false
        }
      ],
      body: {
        mode: "raw",
        raw: JSON.stringify({
          username: "johndoe",
          email: "john@example.com",
          role: "user"
        }),
        options: {
          raw: {
            language: "json"
          }
        }
      }
    },
    response: []
  };


  // Convert request object to internal format for saving state
  function getRequestState() {
    return item
  }

  // Update the state change watcher
  $: {
    const state = {
      ...getRequestState(),
      responseData,
      responseStatus,
      responseHeaders,
      responseTime,
      responseSize
    }
    onStateChange(state)
  }

  function addHeader() {
    if (!item) return;
    item.request.header = [...item.request.header, { key: '', value: '', disabled: false }]
  }

  function addQueryParam() {
    if (!item) return;
    if (!item.request.url.query) item.request.url.query = []
    item.request.url.query = [...item.request.url.query, { key: '', value: '', disabled: false }]
  }

  function addFormDataField() {
    if (!item) return;
    if (!item.request.body) item.request.body = { mode: 'formdata', formdata: [] }
    if (!item.request.body.formdata) item.request.body.formdata = []
    item.request.body.formdata = [...item.request.body.formdata, { key: '', type: 'text', value: '' }]
  }

  // Add focus handlers
  function handleParamKeyFocus(index: number) {
    if (!item) return;
    if (item.request.url.query && index === item.request.url.query.length - 1) {
      addQueryParam()
    }
  }

  function handleHeaderKeyFocus(index: number) {
    if (!item) return;
    if (index === item.request.header.length - 1) {
      addHeader()
    }
  }

  function handleFormDataKeyFocus(index: number) {
    if (!item) return;
    if (item.request.body && item.request.body.formdata && index === item.request.body.formdata.length - 1) {
      addFormDataField()
    }
  }

  // Update URL change handler
  function handleParamChange() {
    const newUrl = buildUrl()
    if (!item) return;
    if (newUrl !== item.request.url.raw) {
      item.request.url.raw = newUrl
    }
    
    // Ensure there's always an empty row at the end
    const query = item.request.url.query || []
    const lastParam = query[query.length - 1]
    if (lastParam && (lastParam.key || lastParam.value)) {
      addQueryParam()
    }
  }

  // Update buildUrl to use request object
  function buildUrl() {
    try {
      if (!item) return '';
      const urlParts = item.request.url.raw.split('?')[0]
      const baseUrl = new URL(urlParts)
      
      baseUrl.search = ''
      
      item.request.url.query?.forEach(param => {
        if (!param.disabled && param.key.trim()) {
          baseUrl.searchParams.set(param.key.trim(), param.value?.trim() || '')
        }
      })
      
      return baseUrl.toString()
    } catch (error) {
      return item?.request?.url?.raw || ''
    }
  }

  // Update onMount to handle PostmanRequest
  onMount(() => {
    console.log('REQUEST TAB ON MOUNT CALLED')
    if (savedState) {
      console.log("SAVED STATE FOUND")
      item = savedState
    } else {
      console.log("NO SAVED STATE")
      item = defaultCollection.items[0]
    }
    
  })

  function handleTabKey(event: KeyboardEvent) {
    if (event.key === 'Tab') {
      event.preventDefault();
      const target = event.target as HTMLTextAreaElement;
      const start = target.selectionStart;
      const end = target.selectionEnd;
      
      const value = target.value;
      target.value = value.substring(0, start) + '\t' + value.substring(end);
      
      target.selectionStart = target.selectionEnd = start + 1;
    }
  }

  async function handleSend() {
    if (item) {
      item.id = tabID
    }


    isLoading = true
    responseData = ''
    responseStatus = ''
    const finalUrl = buildUrl()
    
    // Build headers object
    const headers: { [key: string]: string } = {}
    if (item?.request?.header) {
      item.request.header.forEach(header => {
        if (!header.disabled && header.key.trim() && header.value.trim()) {
          headers[header.key.trim()] = header.value.trim()
        }
      })
    }

    // Handle form data and files
    let formDataObj: { [key: string]: string } = {}
    let files: Array<{fieldName: string, fileName: string, content: string}> = []
    
    if (item?.request?.body?.mode === 'formdata') {
      for (const formItem of item.request.body.formdata || []) {
        if (formItem.key.trim()) {
          if (formItem.type === 'file' && formItem.value instanceof File) {
            // Convert file to base64
            const base64 = await new Promise<string>((resolve) => {
              const reader = new FileReader()
              reader.onloadend = () => {
                if (typeof reader.result === 'string') {
                  const parts = reader.result.split(',')
                  resolve(parts[1] || '')
                } else {
                  resolve('')
                }
              }
              reader.readAsDataURL(formItem.value)
            })
            
            files.push({
              fieldName: formItem.key.trim(),
              fileName: formItem.value.name,
              content: base64
            })
          } else if (formItem.type === 'text' && formItem.value?.trim()) {
            formDataObj[formItem.key.trim()] = formItem.value.trim()
          }
        }
      }
    }
    
    const requestParams = {
      method: item?.request?.method,
      url: finalUrl,
      headers: headers,
      body: item?.request?.body?.mode === 'raw' ? item?.request?.body?.raw : '',
      formData: formDataObj,
      files: files,
      bodyType: item?.request?.body?.mode,
      convertValues: true
    }
    
    MakeRequest(requestParams)
      .then(response => {
        responseData = response.body
        responseHeaders = response.headers
        responseStatus = response.status
        responseTime = response.timeInSeconds
        responseSize = response.responseSize
      })
      .catch(err => {
        console.error(err)
        responseData = `Error: ${err}`
      })
      .then(() => {
        isLoading = false
      })
  }

  function handleKeyPress(event) {
    if (event.key === 'Enter') {
      handleSend()
    }
  }

  // Update toggleParamEnabled to handle URL updates
  function toggleParamEnabled(param) {
    param.disabled = !param.disabled
    handleParamChange()
  }

  function toggleHeaderEnabled(header) {
    header.disabled = !header.disabled;
  }
</script>

<div class="bg-gray-800 rounded-lg p-4">
  <!-- URL Bar -->
  <div class="mb-4">
    <div class="flex flex-col md:flex-row items-center gap-4">
      <select 
        bind:value={item.request.method}
        class="bg-gray-800 text-white font-medium px-3 py-1.5 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 border border-gray-700 appearance-none cursor-pointer [&>option]:bg-gray-800 [&>option]:text-white text-sm"
        style="background-image: url('data:image/svg+xml;utf8,<svg fill=white xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 24 24%22><path d=%22M7 10l5 5 5-5z%22/></svg>'); background-repeat: no-repeat; background-position: right 8px center; padding-right: 32px;"
      >
        {#each methods as method}
          <option value={method}>{method}</option>
        {/each}
      </select>
      <input 
        bind:value={item.request.url.raw}
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
            {#each item.request.url.query || [] as param, i}
              <div class="grid grid-cols-1 md:grid-cols-[24px_2fr_2fr] gap-4 items-end">
                <div class="flex items-end">
                  <div class="h-[24px] flex items-center">
                    <input 
                      type="checkbox"
                      checked={!param.disabled}
                      on:change={() => toggleParamEnabled(param)}
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
                    class="param-input w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm {param.disabled ? 'opacity-50' : ''}"
                    disabled={param.disabled}
                  >
                </div>
                <div>
                  <label class="block text-xs text-gray-400 mb-0.5">Value</label>
                  <input 
                    bind:value={param.value}
                    on:input={handleParamChange}
                    type="text" 
                    class="param-input w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm {param.disabled ? 'opacity-50' : ''}"
                    disabled={param.disabled}
                  >
                </div>
              </div>
            {/each}
          </div>
        {:else if activeTab === 'headers'}
          <div class="space-y-4 max-h-[200px] overflow-y-auto pr-2 pb-2">
            {#each item.request.header as header, i}
              <div class="grid grid-cols-1 md:grid-cols-[24px_2fr_2fr] gap-4 items-end">
                <div class="flex items-end">
                  <div class="h-[24px] flex items-center">
                    <input 
                      type="checkbox"
                      checked={!header.disabled}
                      on:change={() => toggleHeaderEnabled(header)}
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
                    disabled={header.disabled}
                  >
                </div>
                <div>
                  <label class="block text-xs text-gray-400 mb-0.5">Value</label>
                  <input 
                    bind:value={header.value}
                    type="text" 
                    class="header-input w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm {!header.enabled ? 'opacity-50' : ''}"
                    disabled={header.disabled}
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
                  bind:group={item.request.body.mode}
                  value="raw"
                  class="mr-2"
                >
                Raw
              </label>
              <label class="flex items-center">
                <input
                  type="radio"
                  bind:group={item.request.body.mode}
                  value="formdata"
                  class="mr-2"
                >
                Form Data
              </label>
            </div>

            {#if item.request.body.mode === 'raw'}
              <textarea
                bind:value={item.request.body.raw}
                placeholder="Enter request body (JSON, text, etc.)"
                on:keydown={handleTabKey}
                class="w-full h-40 bg-gray-700 text-gray-300 px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
              ></textarea>
            {:else if item.request.body.mode === 'formdata'}
              <div class="space-y-4">
                {#each item.request.body.formdata || [] as formItem, i}
                  <div class="grid grid-cols-1 md:grid-cols-[2fr_2fr_3fr] gap-4 items-end">
                    <div>
                      <label class="block text-xs text-gray-400 mb-1">Key</label>
                      <input 
                        bind:value={formItem.key}
                        type="text"
                        on:focus={() => handleFormDataKeyFocus(i)}
                        class="w-full bg-gray-700 text-gray-300 px-3 py-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                      >
                    </div>
                    <div>
                      <label class="block text-xs text-gray-400 mb-1">Type</label>
                      <select
                        bind:value={formItem.type}
                        class="w-full bg-gray-800 text-white font-medium px-3 py-1.5 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 border border-gray-700 appearance-none cursor-pointer [&>option]:bg-gray-800 [&>option]:text-white text-sm"
                        style="background-image: url('data:image/svg+xml;utf8,<svg fill=white xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 24 24%22><path d=%22M7 10l5 5 5-5z%22/></svg>'); background-repeat: no-repeat; background-position: right 8px center; padding-right: 32px;"
                      >
                        <option value="text">Text</option>
                        <option value="file">File</option>
                      </select>
                    </div>
                    <div>
                      <label class="block text-xs text-gray-400 mb-1">Value</label>
                      {#if formItem.type === 'text'}
                        <input 
                          bind:value={formItem.value}
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
                          on:change={(e: Event) => {
                            const target = e.target as HTMLInputElement;
                            formItem.src = target.files?.[0]?.name || '';
                            // Store the file object temporarily for upload
                            formItem.value = target.files?.[0];
                          }}
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
