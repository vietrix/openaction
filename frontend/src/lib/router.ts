import { writable } from 'svelte/store';

type LocationState = {
  pathname: string;
};

const getLocation = (): LocationState => {
  if (typeof window === 'undefined') {
    return { pathname: '/' };
  }
  return { pathname: window.location.pathname };
};

export const location = writable<LocationState>(getLocation());

const updateLocation = () => {
  location.set(getLocation());
};

if (typeof window !== 'undefined') {
  window.addEventListener('popstate', updateLocation);
}

export const navigate = (to: string) => {
  if (typeof window === 'undefined') return;
  if (to === window.location.pathname) return;
  window.history.pushState({}, '', to);
  updateLocation();
};

export const link = (node: HTMLAnchorElement) => {
  const handleClick = (event: MouseEvent) => {
    if (event.defaultPrevented) return;
    if (event.button !== 0) return;
    if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
    if (node.target && node.target !== '_self') return;

    const href = node.getAttribute('href');
    if (!href || href.startsWith('http') || href.startsWith('mailto:') || href.startsWith('tel:')) {
      return;
    }

    event.preventDefault();
    navigate(href);
  };

  node.addEventListener('click', handleClick);

  return {
    destroy() {
      node.removeEventListener('click', handleClick);
    },
  };
};
