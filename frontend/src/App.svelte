<script lang="ts">
  import logo from './assets/images/logo-universal.png'
  import RequestTab from './components/RequestTab.svelte'
  import { onMount } from 'svelte'
  import { GetCollections, CreateCollection } from '../wailsjs/go/main/App.js'
  import Modal from './components/Modal.svelte'
  import type { PostmanCollection, PostmanItem, Variable } from './types/postman.ts';

  interface Tab {
    id: string;
    request: PostmanItem;
    variable: Variable[];
    title: string;
    hasChanges: boolean; 
  }

  let tabs: Tab[] = [];
  let activeTabId: string | null = null;
  let selectedRequest: PostmanItem | null = null;
  let collections: PostmanCollection[] = []; 
  let currentItem: PostmanItem

  // Add state for expanded collections
  let expandedCollections: Record<string, boolean> = {};

  // Toggle collection expansion
  function toggleCollection(collectionId: string) {
    expandedCollections[collectionId] = !expandedCollections[collectionId];
  }

  function handleKeyDown(event: KeyboardEvent) {
    if (event.ctrlKey && event.key === 's' && activeTabId !== null && currentItem) {
      event.preventDefault();
      console.log('Saving changes...');

      for (const collection of collections) {
        const itemIndex = collection.item.findIndex((item: PostmanItem) => item.id === currentItem.id);
        if (itemIndex !== -1) {
          collection.item[itemIndex] = currentItem;
          localStorage.setItem(collection.info.name, JSON.stringify(collection));
          
          tabs = tabs.map(tab => {
            if (tab.title === currentItem.name) {
              return { ...tab, hasChanges: false };
            }
            return tab;
          });

          console.log(`Saved changes to collection: ${collection.info.name}`);
          return;
        }
      }
      
      // If we get here, the item wasn't found in any collection
      console.warn('Could not find collection for this item');
    }
  }

  // Modified onMount to get tabs from collection
  onMount(() => {
    // Add event listener
    window.addEventListener('keydown', handleKeyDown);

    collections = [];
    for(let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key) {
        const value = localStorage.getItem(key);
        if (value) {
          try {
            const collection = JSON.parse(value);
            if (collection.info && collection.item) {
              collections.push(collection);
            }
          } catch (error) {
            console.error('Invalid collection in localStorage:', key);
          }
        }
      }
    }

    // Clean up event listener
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  })

  // Function to handle request item click
  function handleRequestClick(request: PostmanItem, variable: Variable[]) {
    const newTab: Tab = {
      id: Date.now().toString(),
      request: request,
      title: request.name,
      hasChanges: false,
      variable: variable
    };

    // Check if tab already exists
    const existingTab = tabs.find(tab => tab.title === request.name);
    if (existingTab) {
      activeTabId = existingTab.id;
      selectedRequest = existingTab.request;
      return;
    }

    tabs = [...tabs, newTab];
    activeTabId = newTab.id;
    selectedRequest = request;
  }

  function closeTab(tabId: string, event: MouseEvent) {
    event.stopPropagation();
    const index = tabs.findIndex(tab => tab.id === tabId);
    const currentTabs = tabs.filter(tab => tab.id !== tabId);
    
    tabs = currentTabs;
    
    // If we're closing the active tab
    if (activeTabId === tabId) {
      if (currentTabs.length > 0) {
        // Try to activate the tab at the same index, or the last tab
        const newActiveIndex = Math.min(index, currentTabs.length - 1);
        activeTabId = currentTabs[newActiveIndex].id;
        selectedRequest = currentTabs[newActiveIndex].request;
      } else {
        activeTabId = null;
        selectedRequest = null;
      }
    }
  }

  function onStateChange(state: PostmanItem) {
    console.log("SAVED STATE CALLED", state)
    if (state.id === "default") {
      return
    }

    // Update tabs properly with reactivity
    tabs = tabs.map(tab => {
      if (tab.title === state.name) {
        let updatedTab = { ...tab, hasChanges: true };
        console.log("Updated tab:", updatedTab);
        return updatedTab;
      }
      return tab;
    });

    console.log('TABS ', tabs)
    currentItem = state
    return

    
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
      const collection: PostmanCollection = JSON.parse(text);

      if (!collection.info || !collection.item) {
        throw new Error('Invalid Postman collection format');
      }

      localStorage.setItem(collection.info.name, text);
      target.value = '';
      
      alert('Collection imported successfully!');

      collections = [...collections, collection];
    } catch (error) {
      console.error('Error importing collection:', error);
      alert('Failed to import collection. Please make sure it\'s a valid Postman collection file.');
    }
  }

  // Helper function to get method color
  function getMethodColor(method: string) {
    return method === 'GET' ? 'text-blue-400' : 
           method === 'POST' ? 'text-green-400' : 
           method === 'PUT' ? 'text-yellow-400' : 
           method === 'DELETE' ? 'text-red-400' : 'text-gray-400'
  }
</script>

<main class="min-h-screen bg-gray-900 text-gray-100 text-sm flex" on:keydown={handleKeyDown}>
  <!-- Sidebar -->
  <div class="w-64 bg-gray-800 border-r border-gray-700 flex-shrink-0">
    <div class="p-4">
      <img src={logo} alt="Logo" class="h-8 w-auto mb-6">
      
      <!-- Add Import Button above New Collection -->
      <label class="w-full mb-4 px-4 py-2 bg-gray-600 hover:bg-gray-700 rounded-lg text-white flex items-center justify-center cursor-pointer">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM6.293 6.707a1 1 0 010-1.414l3-3a1 1 0 011.414 0l3 3a 1 1 0 01-1.414 1.414L11 5.414V13a1 1 0 11-2 0V5.414L7.707 6.707a1 1 0 01-1.414 0z" clip-rule="evenodd" />
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

      <!-- Replace Dynamic Collection List with Imported Collection Items -->
      <div class="space-y-4">
        {#each collections as collection}
          <div class="space-y-2">
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div 
              class="flex items-center justify-between p-2 hover:bg-gray-700 rounded-lg cursor-pointer"
              on:click={() => toggleCollection(collection.info._postman_id)}
            >
              <div class="flex items-center">
                <svg 
                  xmlns="http://www.w3.org/2000/svg" 
                  class="h-4 w-4 mr-2 transform transition-transform {expandedCollections[collection.info._postman_id] ? 'rotate-90' : ''}" 
                  viewBox="0 0 20 20" 
                  fill="currentColor"
                >
                  <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a 1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                </svg>
                <span class="text-gray-300 font-medium">{collection.info.name}</span>
              </div>
              <span class="text-xs text-gray-500">{collection.item?.length || 0} requests</span>
            </div>
            
            {#if expandedCollections[collection.info._postman_id]}
              <div class="pl-4 space-y-1 transition-all">
                {#each collection.item || [] as request}
                  <!-- svelte-ignore a11y-click-events-have-key-events -->
                  <div 
                    class="w-full text-left px-2 py-1 rounded hover:bg-gray-700 flex items-center group cursor-pointer {selectedRequest?.id === request.id ? 'bg-gray-700' : ''}"
                    on:click={() => handleRequestClick(request, collection.variable)}
                  >
                    <span class="text-xs font-medium {getMethodColor(request.request.method)} w-12 flex-shrink-0">{request.request.method}</span>
                    <span class="text-gray-400 text-sm truncate flex-1 min-w-0">{request.name}</span>
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {/each}
      </div>

      <!-- Comment out original collections section -->
      <!-- <div class="space-y-4">
        {#each collections as collection}
          <div class="space-y-2">
            <div class="flex items-center">
              <span class="text-gray-300 font-medium">{collection.name}</span>
            </div>
            <div class="pl-4 space-y-1">
              {#each collection.requests as request}
                <div class="w-full text-left px-2 py-1 rounded hover:bg-gray-700 flex items-center">
                  <span class="text-xs font-medium {request.method === 'GET' ? 'text-blue-400' : 
                    request.method === 'POST' ? 'text-green-400' : 
                    request.method === 'PUT' ? 'text-yellow-400' : 
                    request.method === 'DELETE' ? 'text-red-400' : 'text-gray-400'} w-12">{request.method}</span>
                  <span class="text-gray-400 text-sm truncate">{request.name}</span>
                </div>
              {/each}
            </div>
          </div>
        {/each}
      </div> -->
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
      {#if tabs.length > 0}
        <div class="flex items-center gap-1 mb-4 overflow-x-auto border-b border-gray-700">
          {#each tabs as tab}
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div 
              class="flex-shrink-0 flex items-center group cursor-pointer {activeTabId === tab.id ? 'border-b-2 border-blue-500' : ''}"
              on:click={() => {
                activeTabId = tab.id;
                selectedRequest = tab.request;
              }}
            >
              <div class="px-3 py-2 flex items-center gap-2 hover:bg-gray-800 rounded-t">
                <span class="text-xs font-medium {getMethodColor(tab.request.request.method)} flex-shrink-0">{tab.request.request.method}</span>
                <div class="flex items-center gap-2">
                  <span class="text-sm text-gray-300 truncate max-w-[150px]">{tab.title}</span>
                  {#if tab.hasChanges}
                    <div class="w-1.5 h-1.5 rounded-full bg-blue-500"></div>
                  {/if}
                </div>
                <button
                  class="ml-2 p-1 rounded-full hover:bg-gray-700 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0"
                  on:click={(e) => closeTab(tab.id, e)}
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </div>
          {/each}
        </div>

        <!-- Active Tab Content -->
        {#if activeTabId}
          {#each tabs as tab}
            {#if activeTabId === tab.id}
              <RequestTab 
                tabID={tab.request.id}
                savedState={tab.request}
                variable={tab.variable}
                onStateChange={onStateChange}
              />
            {/if}
          {/each}
        {/if}
      {:else}
        <!-- Show empty state when no tabs are open -->
        <div class="flex items-center justify-center h-[calc(100vh-200px)]">
          <div class="text-center max-w-md">
            <div class="bg-gray-800 p-8 rounded-lg shadow-lg">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto mb-6 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              <h3 class="text-xl font-medium text-gray-100 mb-2">No Request Selected</h3>
              <p class="text-gray-400 mb-6">Select a request from the sidebar to start testing your API endpoints.</p>
              <div class="flex justify-center gap-4">
                <button 
                  class="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg text-sm text-white flex items-center gap-2"
                  on:click={() => showNewCollectionModal = true}
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd" />
                  </svg>
                  New Collection
                </button>
                <label class="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-lg text-sm text-white flex items-center gap-2 cursor-pointer">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
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
              </div>
            </div>
          </div>
        </div>
      {/if}
    </div>
  </div>
</main>
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
  :global(.tab-container) {
    max-width: 100%;
    overflow-x: auto;
    scrollbar-width: thin;
    scrollbar-color: theme('colors.gray.600') theme('colors.gray.800');
  }
  
  :global(.tab-container::-webkit-scrollbar) {
    height: 6px;
  }
  
  :global(.tab-container::-webkit-scrollbar-track) {
    background: theme('colors.gray.800');
  }
  
  :global(.tab-container::-webkit-scrollbar-thumb) {
    background-color: theme('colors.gray.600');
    border-radius: 3px;
  }
</style>
