import { type Ref, onMounted, onUnmounted, ref } from 'vue';
import { type GameState, defaultGameState, parseGame } from './state';

interface Endpoints {
  PUT: {
    '/api/join': {
      request: { name: string };
      response: { playerIndex: number };
    };
    '/api/choose': {
      request: { player: number; white: number; black: number };
    };
    '/api/reset': {};
    '/api/vote': {
      request: { player: number; fighter: number };
    };
  };
  PATCH: {
    '/api/game': {
      request: { Goal?: number; HandSize?: number };
    };
  };
}

type ApiRequest<T> = T extends { request: infer U } ? U : undefined;
type ApiResponse<T> = T extends { response: infer U } ? U : undefined;

export async function callApi<
  M extends keyof Endpoints,
  P extends keyof Endpoints[M] & string,
>(
  method: M,
  endpoint: P,
  body?: ApiRequest<Endpoints[M][P]>,
): Promise<ApiResponse<Endpoints[M][P]>> {
  const res = await fetch(endpoint, {
    method,
    body: body
      ? new URLSearchParams(body as Record<string, string>)
      : undefined,
  });
  return res.status === 200
    ? res.json()
    : (undefined as ApiResponse<Endpoints[M][P]>);
}

export function useServerGameState(
  onReset?: (kind: string) => void,
): Ref<GameState> {
  let source: EventSource;
  let controller: AbortController;
  const gamestate = ref(defaultGameState);

  onMounted(() => {
    source = new EventSource('/api/events');
    controller = new AbortController();
    const signal = controller.signal;

    function onEvent(
      event: string,
      listener: (event: MessageEvent) => void,
    ): void {
      source.addEventListener(event, listener, { signal });
    }

    onEvent('gameupdate', event => (gamestate.value = parseGame(event.data)));

    if (onReset) {
      onEvent('reset', event => onReset(event.data));
    }
  });

  onUnmounted(() => {
    controller.abort();
    source.close();
  });

  return gamestate;
}
