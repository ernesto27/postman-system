<script lang="ts">
  import logo from './assets/images/logo-universal.png'
  import RequestTab from './components/RequestTab.svelte'
  import { onMount } from 'svelte'
  import { GetCollections, CreateCollection } from '../wailsjs/go/main/App.js'
  import Modal from './components/Modal.svelte'
  import type { PostmanCollection, PostmanItem } from './types/postman.ts';

  interface Tab {
    id?: number;
    title?: string;
  }

  let tabs: Tab[] = []
  let activeTabId = 1
  let nextTabId = 2
  let collections = [] // Initialize as empty array

  // Modified onMount to get tabs from collection
  onMount(() => {
    GetCollections()
      .then(response => {
        collections = Array.isArray(response) ? response : []
        console.log("Collections loaded:", collections)
        
        // Get collection from localStorage
        const collection = localStorage.getItem("collection_1")
        if (collection) {
          const collObj = JSON.parse(collection)
          // Create tabs from collection items
          tabs = collObj.item.map((item, index) => ({
            id: index,
            title: item.name || `Request ${index + 1}`
          }))
          nextTabId = Math.max(...tabs.map(tab => tab.id ?? 0)) + 1
          console.log("TABS" , tabs)
        } else {
          tabs = [{ id: 1, title: 'New Request 1' }]
        }
        activeTabId = tabs[0]?.id || 1
      })
      .catch(err => {
        console.error("Failed to load collections:", err)
        collections = []
        tabs = [{ id: 1, title: 'New Request 1' }]
      })
  })

  function addNewTab() {
    const newTab = {
      id: nextTabId,
      title: `New Request ${nextTabId}`
    }
    tabs = [...tabs, newTab]
    activeTabId = nextTabId
    nextTabId++
  }

  function removeTab(tabId: number | undefined) {
    tabs = tabs.filter(tab => tab.id !== tabId)
    if (activeTabId === tabId && tabs.length > 0) {
      activeTabId = tabs[tabs.length - 1].id ?? 1
    }

    // Remove from collection in localStorage
    const collection = localStorage.getItem("collection_1")
    if (collection) {
      const collObj = JSON.parse(collection)
      collObj.item = collObj.item.filter(item => item.id !== tabId.toString())
      localStorage.setItem('collection_1', JSON.stringify(collObj))
    }
  }

  function saveTabState(state: PostmanItem) {
    console.log("SAVED STATE CALLED", state)
    if (state.id === "default") {
      return
    }

    const collection = localStorage.getItem("collection_1")
    if (!collection) {
      const collectionToSave: PostmanCollection = {
        id: state.id,
        name: "API Testing Collection",
        item: [state]  // Changed from items to item
      }
  
      localStorage.setItem('collection_1', JSON.stringify(collectionToSave))
      return
    } else {
      const collObj = JSON.parse(collection)
      const itemIndex = collObj.item.findIndex((item: PostmanItem) => item.id === state.id)  // Changed from items to item
      if (itemIndex === -1) {
        collObj.item.push(state)  // Changed from items to item
      } else {
        collObj.item[itemIndex] = state  // Changed from items to item
      }

      localStorage.setItem('collection_1', JSON.stringify(collObj))
    }
  }

  function getTabState(id: number) {
    const collection = localStorage.getItem("collection_1")
    const collObj = collection ? JSON.parse(collection) : null
    console.log(collObj)

    if (!collObj) {
      return null
    }
    
    const resp = collObj.item.find((item: PostmanItem) => item.id === id.toString()) || null  // Changed from items to item
    console.log("GET TAB STATE CALLED", resp)
    return resp;
  }

  let isSidebarOpen = true;
  let showNewCollectionModal = false;
  let newCollectionName = '';

  function createNewCollection() {
    if (newCollectionName.trim()) {
      CreateCollection(newCollectionName.trim(), "1")
        .then(response => {
          collections = [...collections, { name: newCollectionName.trim(), requests: [] }];
          newCollectionName = '';
          showNewCollectionModal = false;
          console.log("New collection created:", response)
        })
        .catch(err => {
          alert("Failed to create collection. Please try again.")
          console.error("Failed to create collection:", err)
        })
    }
  }

  let isHistoryOpen = false;
  const mockHistory = [
    { method: 'GET', url: 'https://api.example.com/users', timestamp: '2 mins ago' },
    { method: 'POST', url: 'https://api.example.com/orders', timestamp: '15 mins ago' },
    { method: 'PUT', url: 'https://api.example.com/products/1', timestamp: '1 hour ago' },
    { method: 'DELETE', url: 'https://api.example.com/posts/5', timestamp: '2 hours ago' }
  ];

  // Add new function to handle file import
  async function handleFileImport(event: Event) {
    const target = event.target as HTMLInputElement;
    const file = target.files?.[0];
    if (!file) return;

    try {
      const text = await file.text();
      const collection = JSON.parse(text);

      // Basic validation that it's a Postman collection
      if (!collection.info || !collection.item) {
        throw new Error('Invalid Postman collection format');
      }

      // Store in localStorage
      localStorage.setItem('collection_1', text);

      
      // Reset file input
      target.value = '';
      
      alert('Collection imported successfully!');
    } catch (error) {
      console.error('Error importing collection:', error);
      alert('Failed to import collection. Please make sure it\'s a valid Postman collection file.');
    }
  }
</script>

<main class="min-h-screen bg-gray-900 text-gray-100 text-sm flex">
  <!-- Sidebar -->
  <div class="w-64 bg-gray-800 border-r border-gray-700 flex-shrink-0">
    <div class="p-4">
      <img src={logo} alt="Logo" class="h-8 w-auto mb-6">
      
      <!-- Add Import Button above New Collection -->
      <label class="w-full mb-4 px-4 py-2 bg-gray-600 hover:bg-gray-700 rounded-lg text-white flex items-center justify-center cursor-pointer">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM6.293 6.707a1 1 0 010-1.414l3-3a1 1 0 011.414 0l3 3a1 1 0 01-1.414 1.414L11 5.414V13a1 1 0 11-2 0V5.414L7.707 6.707a1 1 0 01-1.414 0z" clip-rule="evenodd" />
        </svg>
        Import Collection
        <input 
          type="file" 
          accept=".json"
          class="hidden" 
          on:change={handleFileImport}
        >
      </label>

      <!-- Add New Collection Button -->
      <button
        class="w-full mb-4 px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg text-white flex items-center justify-center"
        on:click={() => showNewCollectionModal = true}
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd" />
        </svg>
        New Collection
      </button>

      <!-- History Button -->
      <button
        class="w-full mb-4 px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg text-white flex items-center justify-between"
        on:click={() => isHistoryOpen = !isHistoryOpen}
      >
        <div class="flex items-center">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
          </svg>
          History
        </div>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 transform {isHistoryOpen ? 'rotate-180' : ''}" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
        </svg>
      </button>

      <!-- History List -->
      {#if isHistoryOpen}
        <div class="mb-6 space-y-2">
          {#each mockHistory as item}
            <div class="p-2 hover:bg-gray-700 rounded-lg cursor-pointer">
              <div class="flex items-center gap-2">
                <span class="text-xs font-medium {
                  item.method === 'GET' ? 'text-blue-400' : 
                  item.method === 'POST' ? 'text-green-400' : 
                  item.method === 'PUT' ? 'text-yellow-400' : 
                  'text-red-400'
                }">{item.method}</span>
                <span class="text-gray-400 text-xs">{item.timestamp}</span>
              </div>
              <div class="text-gray-300 text-sm truncate mt-1">
                {item.url}
              </div>
            </div>
          {/each}
        </div>
      {/if}

      <!-- Dynamic Collection List -->
      <div class="space-y-4">
        {#each collections as collection}
          <div class="space-y-2">
            <div class="flex items-center">
              <span class="text-gray-300 font-medium">{collection.name}</span>
            </div>
            <div class="pl-4 space-y-1">
              <!-- {#each collection.requests as request}
                <div class="w-full text-left px-2 py-1 rounded hover:bg-gray-700 flex items-center">
                  <span class="text-xs font-medium {request.method === 'GET' ? 'text-blue-400' : 
                    request.method === 'POST' ? 'text-green-400' : 
                    request.method === 'PUT' ? 'text-yellow-400' : 
                    request.method === 'DELETE' ? 'text-red-400' : 'text-gray-400'} w-12">{request.method}</span>
                  <span class="text-gray-400 text-sm truncate">{request.name}</span>
                </div>
              {/each} -->
            </div>
          </div>
        {/each}
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
              on:click={() => activeTabId = tab.id ?? 1}
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
            tabID={tab.id.toString()}
            savedState={getTabState(tab.id)}
            onStateChange={(state) => saveTabState(state)}
          />
        {/if}
      {/each}
    </div>
  </div>
</main>

<!-- New Collection Modal -->
<Modal show={showNewCollectionModal} title="Create New Collection">
  <div class="space-y-4">
    <input
      type="text"
      bind:value={newCollectionName}
      placeholder="Collection Name"
      class="w-full px-3 py-2 bg-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
    <div class="flex justify-end gap-2">
      <button
        class="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg"
        on:click={() => showNewCollectionModal = false}
      >
        Cancel
      </button>
      <button
        class="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg"
        on:click={createNewCollection}
      >
        Create
      </button>
    </div>
  </div>
</Modal>

<style>
  /* Add any additional styles here */
</style>
