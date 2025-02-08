<script lang="ts">
  import type { PostmanItem, Variable } from '../types/postman';

  export let item: PostmanItem;
  export let level: number = 0;
  export let variables: Variable[] = [];
  export let selectedRequest: PostmanItem | null;
  export let expandedCollections: Record<string, boolean>;
  export let handleRequestClick: (request: PostmanItem, variables: Variable[]) => void;

  // Helper function to get all child folder IDs
  function getAllChildFolderIds(item: PostmanItem): string[] {
    let ids: string[] = [];
    if (item.item) {
      ids.push(item.id);
      item.item.forEach(child => {
        ids = [...ids, ...getAllChildFolderIds(child)];
      });
    }
    return ids;
  }

  // Modified toggleCollection to handle children
  function toggleCollection(id: string) {
    if (!expandedCollections[id]) {
      // Opening folder - just set its state
      expandedCollections[id] = true;
    } else {
      // Closing folder - remove all child states
      const childIds = getAllChildFolderIds(item);
      childIds.forEach(childId => {
        delete expandedCollections[childId];
      });
      expandedCollections[id] = false;
    }
    expandedCollections = expandedCollections; // Trigger reactivity
  }

  function getMethodColor(method: string) {
    return method === 'GET' ? 'text-blue-400' : 
           method === 'POST' ? 'text-green-400' : 
           method === 'PUT' ? 'text-yellow-400' : 
           method === 'DELETE' ? 'text-red-400' : 'text-gray-400';
  }
</script>

<div class="w-full" style="padding-left: {level * 12}px">
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div 
    class="w-full text-left px-2 py-1 rounded hover:bg-gray-700 flex items-center group cursor-pointer {selectedRequest?.id === item.id ? 'bg-gray-700' : ''}"
    on:click={() => {
      if (item.item) {
        toggleCollection(item.id);
      } else {
        handleRequestClick(item, variables);
      }
    }}
  >
    {#if item.item}
      <div class="flex items-center w-full">
        <svg 
          xmlns="http://www.w3.org/2000/svg" 
          class="h-4 w-4 mr-2 transform transition-transform {expandedCollections[item.id] ? 'rotate-90' : ''}" 
          viewBox="0 0 20 20" 
          fill="currentColor"
        >
          <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
        </svg>
        <span class="text-gray-400 text-sm truncate flex-1 min-w-0">{item.name}</span>
      </div>
    {:else}
      <span class="text-xs font-medium {getMethodColor(item.request.method)} w-12 flex-shrink-0">{item.request.method}</span>
      <span class="text-gray-400 text-sm truncate flex-1 min-w-0">{item.name}</span>
    {/if}
  </div>

  {#if item.item && expandedCollections[item.id]}
    <div class="space-y-1">
      {#each item.item as subItem}
        <svelte:self
          item={subItem}
          level={level + 1}
          variables={variables}
          {selectedRequest}
          {expandedCollections}
          {handleRequestClick}
        />
      {/each}
    </div>
  {/if}
</div>
