<script>
  export let show = false;
  export let title = '';
  
  function closeModal() {
    show = false;
  }
</script>

{#if show}
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-gray-800 rounded-lg p-6 w-96">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-lg font-semibold">{title}</h2>
      </div>
      <slot></slot>
    </div>
  </div>
{/if}
