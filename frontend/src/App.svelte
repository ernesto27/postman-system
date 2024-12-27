<script>
  import logo from './assets/images/logo-universal.png'
  import {Greet} from '../wailsjs/go/main/App.js'

  let resultText = "Please enter your name below 👇"
  let name

  function greet() {
    Greet(name).then(result => resultText = result)
  }
</script>

<main class="min-h-screen bg-gray-900 text-gray-100">
  <div class="max-w-7xl mx-auto p-6">
    <!-- Header -->
    <div class="mb-8">
      <img src={logo} alt="Logo" class="h-8 w-auto mb-6">
      <div class="bg-gray-800 rounded-lg p-4">
        <div class="flex flex-col md:flex-row items-center gap-4">
          <select class="bg-gray-700 text-gray-300 px-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
            <option>GET</option>
            <option>POST</option>
            <option>PUT</option>
            <option>DELETE</option>
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

    <!-- Request/Response Tabs -->
    <div class="bg-gray-800 rounded-lg p-4">
      <div class="border-b border-gray-700 mb-4">
        <div class="flex gap-4">
          <button class="px-4 py-2 text-blue-500 border-b-2 border-blue-500">Params</button>
          <button class="px-4 py-2 text-gray-400 hover:text-gray-300">Headers</button>
          <button class="px-4 py-2 text-gray-400 hover:text-gray-300">Body</button>
        </div>
      </div>
      
      <!-- Content Area -->
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
  </div>
</main>
