import { type Ref, onUnmounted, ref } from 'vue';
import { type GameSettings, type GameState, defaultGameState } from './state';

interface Endpoints {
  GET: {
    '/api/game': {
      response: GameState;
    };
  };
  PATCH: {
    '/api/game': {
      request: Partial<GameSettings>;
    };
  };
  PUT: {
    '/api/join': {
      request: { name: string };
      response: { playerIndex: number };
    };
    '/api/choose': {
      request: { player: number; white: number; black: number[] };
    };
    '/api/reset': {};
    '/api/vote': {
      request: { player: number; fighter: number };
    };
  };
}

type ApiRequest<T> = T extends { request: infer U } ? U : undefined;
type ApiResponse<T> = T extends { response: infer U } ? U : undefined;

function parseBody(body: unknown): URLSearchParams | undefined {
  if (!body) {
    return;
  }
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(body)) {
    if (!Array.isArray(value)) {
      params.set(key, value);
      continue;
    }
    for (const val of value) {
      params.append(key, val);
    }
  }
  return params;
}

function subscribe(
  source: EventSource,
  signal: AbortSignal,
  listeners: { [name: string]: (event: MessageEvent) => void },
) {
  const opts = { signal };
  for (const [name, listener] of Object.entries(listeners)) {
    source.addEventListener(name, listener, opts);
  }
}

export interface UseApi {
  gamestate: Ref<GameState>;
  callApi: <M extends keyof Endpoints, P extends keyof Endpoints[M] & string>(
    method: M,
    endpoint: P,
    body?: ApiRequest<Endpoints[M][P]>,
  ) => Promise<ApiResponse<Endpoints[M][P]>>;
}

export function useApi(onReset?: (kind: string) => void): UseApi {
  const source = new EventSource('/api/events');
  const controller = new AbortController();
  const { signal } = controller;
  const gamestate = ref(defaultGameState);

  onUnmounted(() => {
    controller.abort();
    source.close();
  });

  subscribe(source, signal, {
    gameupdate(event) {
      gamestate.value = JSON.parse(event.data);
    },
    reset(event) {
      onReset?.(event.data);
    },
    shutdown() {
      source.close();
    },
  });

  async function callApi<
    M extends keyof Endpoints,
    P extends keyof Endpoints[M] & string,
  >(
    method: M,
    endpoint: P,
    body?: ApiRequest<Endpoints[M][P]>,
  ): Promise<ApiResponse<Endpoints[M][P]>> {
    const res = await fetch(endpoint, {
      method,
      body: parseBody(body),
      signal,
    });
    return res.status === 200
      ? res.json()
      : (undefined as ApiResponse<Endpoints[M][P]>);
  }

  callApi('GET', '/api/game', undefined)
    .then(state => (gamestate.value = state))
    .catch(console.error);

  return { gamestate, callApi };
}
