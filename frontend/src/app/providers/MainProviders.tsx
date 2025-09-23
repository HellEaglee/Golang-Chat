import { router } from "@/app/routers";
import { store } from "@/app/store";
import type { ReactNode } from "react";
import { Provider } from "react-redux";
import { RouterProvider } from "react-router-dom";

type Props = {
  children?: ReactNode;
};

export const MainProviders = ({ children }: Props) => {
  return (
    <Provider store={store}>
      <RouterProvider router={router} />
      {children}
    </Provider>
  );
};
