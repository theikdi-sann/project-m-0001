import { render, screen, act } from "@testing-library/react";
import { CartProvider, useCart } from "@/context/CartContext";
import { MenuItem } from "@/types/menu";

// Test Component to consume context
const TestComponent = () => {
  const { addToCart, totalItems, totalPrice, cart } = useCart();

  const item: MenuItem = {
    ID: "1",
    CategoryID: "cat1",
    Name: "Burger",
    Price: "10.00",
    IsAvailable: true,
  };

  return (
    <div>
      <div data-testid="count">{totalItems}</div>
      <div data-testid="total">{totalPrice}</div>
      <button onClick={() => addToCart(item)}>Add</button>
    </div>
  );
};

describe("CartContext", () => {
  it("adds items and calculates total", () => {
    render(
      <CartProvider>
        <TestComponent />
      </CartProvider>
    );

    expect(screen.getByTestId("count")).toHaveTextContent("0");
    
    act(() => {
      screen.getByText("Add").click();
    });

    expect(screen.getByTestId("count")).toHaveTextContent("1");
    expect(screen.getByTestId("total")).toHaveTextContent("10");

    act(() => {
      screen.getByText("Add").click();
    });

    expect(screen.getByTestId("count")).toHaveTextContent("2");
    expect(screen.getByTestId("total")).toHaveTextContent("20");
  });
});
