import {Body, Controller, Delete, Get, Param, Post, Put} from '@nestjs/common';
import {Category} from './category.entity';
import {CategoryService} from './category.service';

@Controller('categories')
export class CategoryController {

  constructor(
    private readonly categoryService: CategoryService,
  ) {
  }

  @Get()
  public async getCategories(): Promise<Category[]> {
    return this.categoryService.findAllCategories();
  }

  @Post()
  public async createCategory(@Body()category: Category): Promise<any> {
    return this.categoryService.createCategory(category);
  }

  @Put()
  public async editCategory(@Body()category: Category): Promise<any> {
    // return this.categoryService.editCategory(category);
  }

  @Delete(':id')
  public async deleteCategory(@Param('id')id: number): Promise<Category> {
    return this.categoryService.deleteCategory(id);
  }
}
