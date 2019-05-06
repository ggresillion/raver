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

  @Put(':id')
  public async editCategory(@Param('id')id: number, @Body()category: Category): Promise<any> {
    return this.categoryService.editCategory(id, category);
  }

  @Delete(':id')
  public async deleteCategory(@Param('id')id: number): Promise<Category> {
    return this.categoryService.deleteCategory(id);
  }
}
