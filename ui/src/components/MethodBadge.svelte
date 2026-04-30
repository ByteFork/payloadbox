<script lang="ts">
import { methodInfo } from "../lib/format";
export let method: string = "";
export let small: boolean = false;
$: info = methodInfo(method);
</script>

<span
  class="inline-block text-center whitespace-nowrap font-mono font-semibold"
  style="color: {info.color};
         background: {info.bg};
         font-size: {small ? 9 : 10}px;
         letter-spacing: 0.04em;
         padding: {small ? '2px 5px' : '2px 6px'};
         border-radius: 4px;
         min-width: {small ? 38 : 50}px;"
>{method.toUpperCase()}</span>
