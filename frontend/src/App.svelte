<script>
  import logo from './assets/images/logo-universal.png'
  import {Greet} from '../wailsjs/go/main/App.js'

  let resultText = "Please enter your name below 👇"
  let name
  let selectedMethod = 'GET'
  const methods = ['GET', 'POST', 'PUT', 'DELETE']

  function greet() {
    Greet(name).then(result => resultText = result)
  }
</script>

<main class="min-h-screen bg-gray-900 text-gray-100">
  <div class="max-w-7xl mx-auto p-6">
    <!-- Header with Logo and URL Bar -->
    <div class="mb-8">
      <img src={logo} alt="Logo" class="h-8 w-auto mb-6">
      <div class="bg-gray-800 rounded-lg p-4">
        <div class="flex flex-col md:flex-row items-center gap-4">
          <select 
            bind:value={selectedMethod}
            class="bg-gray-800 text-white font-medium px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 border border-gray-700 appearance-none cursor-pointer [&>option]:bg-gray-800 [&>option]:text-white"
            style="background-image: url('data:image/svg+xml;utf8,<svg fill=white xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 24 24%22><path d=%22M7 10l5 5 5-5z%22/></svg>'); background-repeat: no-repeat; background-position: right 8px center; padding-right: 32px;"
          >
            {#each methods as method}
              <option value={method}>{method}</option>
            {/each}
          </select>
          <input 
            type="text" 
            placeholder="Enter request URL"
            class="flex-1 bg-gray-700 text-gray-300 px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
          <button class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg transition-colors">
            Send
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
            <button class="px-4 py-2 text-blue-500 border-b-2 border-blue-500">Params</button>
            <button class="px-4 py-2 text-gray-400 hover:text-gray-300">Headers</button>
            <button class="px-4 py-2 text-gray-400 hover:text-gray-300">Body</button>
          </div>
        </div>
        
        <!-- Request Content Area -->
        <div class="bg-gray-900 rounded-lg p-4">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="col-span-1">
              <label class="block text-sm text-gray-400 mb-2">Key</label>
              <input 
                type="text" 
                class="w-full bg-gray-700 text-gray-300 px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
            </div>
            <div class="col-span-1">
              <label class="block text-sm text-gray-400 mb-2">Value</label>
              <input 
                type="text" 
                class="w-full bg-gray-700 text-gray-300 px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
            </div>
            <div class="col-span-1">
              <label class="block text-sm text-gray-400 mb-2">Description</label>
              <input 
                type="text" 
                class="w-full bg-gray-700 text-gray-300 px-3 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
            </div>
          </div>
        </div>
      </div>

      <!-- Response Section -->
      <div class="bg-gray-800 rounded-lg p-4">
        <div class="border-b border-gray-700 mb-4">
          <div class="flex gap-4">
            <button class="px-4 py-2 text-blue-500 border-b-2 border-blue-500">Response</button>
            <button class="px-4 py-2 text-gray-400 hover:text-gray-300">Headers</button>
          </div>
        </div>
        
        <div class="bg-gray-900 rounded-lg p-4">
          <pre class="font-mono text-sm text-gray-300 whitespace-pre-wrap min-h-[200px]">
{@html `{
  "status": "success",
  "message": "Data retrieved successfully",
  "data": {
    "id": 1,
    "name": "Example"
  }
}`}
          </pre>
        </div>
      </div>
    </div>
  </div>
</main>
