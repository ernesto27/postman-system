<script>
  import logo from './assets/images/logo-universal.png'
  import RequestTab from './components/RequestTab.svelte'

  let tabs = [{ id: 1, title: 'New Request 1' }]
  let activeTabId = 1
  let nextTabId = 2

  // Store tab states in a map
  let tabStates = new Map()

  function addNewTab() {
    const newTab = {
      id: nextTabId,
      title: `New Request ${nextTabId}`
    }
    tabs = [...tabs, newTab]
    activeTabId = nextTabId
    nextTabId++
  }

  function removeTab(tabId) {
    tabs = tabs.filter(tab => tab.id !== tabId)
    // Clean up the state when removing a tab
    tabStates.delete(tabId)
    if (activeTabId === tabId && tabs.length > 0) {
      activeTabId = tabs[tabs.length - 1].id
    }
  }

  function saveTabState(tabId, state) {
    tabStates.set(tabId, state)
  }

  function getTabState(tabId) {
    return tabStates.get(tabId)
  }
</script>

<main class="min-h-screen bg-gray-900 text-gray-100 text-sm">
  <div class="max-w-7xl mx-auto p-6">
    <!-- Header with Logo -->
    <div class="mb-8">
      <img src={logo} alt="Logo" class="h-8 w-auto mb-6">
    </div>

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
</main>

<style>
  /* Add any additional styles here */
</style>
