<script>
  import logo from './assets/images/logo-universal.png'
  import RequestTab from './components/RequestTab.svelte'
  import { onMount } from 'svelte'
  let tabs = []
  let activeTabId = 1
  let nextTabId = 2
  let tabStates = new Map()

  // Load saved state from localStorage
  onMount(() => {
    const savedTabs = localStorage.getItem('postman_tabs')
    const savedStates = localStorage.getItem('postman_states')
    
    if (savedTabs) {
      tabs = JSON.parse(savedTabs)
      nextTabId = Math.max(...tabs.map(tab => tab.id)) + 1
    } else {
      tabs = [{ id: 1, title: 'New Request 1' }]
    }

    if (savedStates) {
      tabStates = new Map(JSON.parse(savedStates))
    }

    activeTabId = tabs[0]?.id || 1
  })

  function addNewTab() {
    const newTab = {
      id: nextTabId,
      title: `New Request ${nextTabId}`
    }
    tabs = [...tabs, newTab]
    activeTabId = nextTabId
    nextTabId++
    saveTabs()
  }

  function removeTab(tabId) {
    tabs = tabs.filter(tab => tab.id !== tabId)
    // Clean up the state when removing a tab
    tabStates.delete(tabId)
    if (activeTabId === tabId && tabs.length > 0) {
      activeTabId = tabs[tabs.length - 1].id
    }
    saveTabs()
  }

  function saveTabState(tabId, state) {
    tabStates.set(tabId, state)
    saveStates()
  }

  function getTabState(tabId) {
    return tabStates.get(tabId)
  }

  // Helper functions to save to localStorage
  function saveTabs() {
    localStorage.setItem('postman_tabs', JSON.stringify(tabs))
  }

  function saveStates() {
    localStorage.setItem('postman_states', JSON.stringify([...tabStates]))
  }

  // Hardcoded collections data
  const collections = [
    {
      id: 1,
      name: 'My Collection 1',
      requests: [
        { id: 1, name: 'Get Users', method: 'GET' },
        { id: 2, name: 'Create User', method: 'POST' }
      ]
    },
    {
      id: 2,
      name: 'API Tests',
      requests: [
        { id: 3, name: 'Authentication', method: 'POST' },
        { id: 4, name: 'List Items', method: 'GET' },
        { id: 5, name: 'Update Item', method: 'PUT' }
      ]
    }
  ]

  let isSidebarOpen = true;
</script>

<main class="min-h-screen bg-gray-900 text-gray-100 text-sm flex">
  <!-- Sidebar -->
  <div class="w-64 bg-gray-800 border-r border-gray-700 flex-shrink-0">
    <div class="p-4">
      <img src={logo} alt="Logo" class="h-8 w-auto mb-6">
      <!-- Static Collection List -->
      <div class="space-y-4">
        <div class="space-y-2">
          <div class="flex items-center">
            <span class="text-gray-300 font-medium">Collection 1</span>
          </div>
          <div class="pl-4 space-y-1">
            <div class="w-full text-left px-2 py-1 rounded hover:bg-gray-700 flex items-center">
              <span class="text-xs font-medium text-blue-400 w-12">GET</span>
              <span class="text-gray-400 text-sm truncate">Get Users</span>
            </div>
            <div class="w-full text-left px-2 py-1 rounded hover:bg-gray-700 flex items-center">
              <span class="text-xs font-medium text-green-400 w-12">POST</span>
              <span class="text-gray-400 text-sm truncate">Create User</span>
            </div>
          </div>
        </div>
        <div class="space-y-2">
          <div class="flex items-center">
            <span class="text-gray-300 font-medium">Collection 2</span>
          </div>
          <div class="pl-4 space-y-1">
            <div class="w-full text-left px-2 py-1 rounded hover:bg-gray-700 flex items-center">
              <span class="text-xs font-medium text-yellow-400 w-12">PUT</span>
              <span class="text-gray-400 text-sm truncate">Update Item</span>
            </div>
            <div class="w-full text-left px-2 py-1 rounded hover:bg-gray-700 flex items-center">
              <span class="text-xs font-medium text-red-400 w-12">DELETE</span>
              <span class="text-gray-400 text-sm truncate">Delete Item</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Main Content -->
  <div class="flex-1">
    <div class="p-6">
      <!-- Toggle Sidebar Button -->
      <button 
        class="mb-4 p-2 hover:bg-gray-800 rounded"
        on:click={() => isSidebarOpen = !isSidebarOpen}
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
        </svg>
      </button>

      <!-- Tabs Bar -->
      <div class="flex items-center gap-2 mb-4 overflow-x-auto">
        {#each tabs as tab}
          <div class="flex items-center">
            <button
              class="px-4 py-2 rounded-t-lg {activeTabId === tab.id ? 'bg-gray-800 text-white' : 'bg-gray-700 text-gray-300 hover:bg-gray-600'}"
              on:click={() => activeTabId = tab.id}
            >
              {tab.title}
            </button>
            {#if tabs.length > 1}
              <button
                class="ml-2 text-gray-500 hover:text-gray-300"
                on:click={() => removeTab(tab.id)}
              >
                ×
              </button>
            {/if}
          </div>
        {/each}
        <button
          class="px-3 py-2 rounded-lg bg-gray-700 hover:bg-gray-600 text-gray-300"
          on:click={addNewTab}
        >
          +
        </button>
      </div>

      <!-- Active Tab Content -->
      {#each tabs as tab}
        {#if activeTabId === tab.id}
          <RequestTab 
            title={tab.title}
            savedState={getTabState(tab.id)}
            onStateChange={(state) => saveTabState(tab.id, state)}
          />
        {/if}
      {/each}
    </div>
  </div>
</main>

<style>
  /* Add any additional styles here */
</style>
